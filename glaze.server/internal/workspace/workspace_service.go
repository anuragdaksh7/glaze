package workspace

import (
	"context"
	"errors"
	"fmt"
	"glaze/config"
	projectDto "glaze/dto/project"
	userDto "glaze/dto/user"
	workspaceDto "glaze/dto/workspace"
	workspaceMemberDto "glaze/dto/workspacemember"
	"glaze/logger"
	"glaze/models"
	"glaze/utils"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type service struct {
	timeout time.Duration
	DB      *gorm.DB
}

func NewService() Service {
	return &service{
		time.Duration(20) * time.Second,
		config.DB,
	}
}

func (s *service) GetAllWorkspaces(c context.Context, userID uuid.UUID) (*workspaceDto.WorkspaceList, error) {
	var result []workspaceDto.WorkspaceWithCounts

	memberSub := s.DB.
		Table("workspace_members").
		Select("workspace_id, COUNT(*) as member_count").
		Group("workspace_id")

	projectSub := s.DB.
		Table("projects").
		Select("workspace_id, COUNT(*) as project_count").
		Group("workspace_id")

	err := s.DB.
		Table("workspaces w").
		Select(`
            w.id,
            w.name,
            w.slug,
            w.billing_plan,
            COALESCE(m.member_count, 0) as member_count,
            COALESCE(p.project_count, 0) as project_count
        `).
		Joins("JOIN workspace_members wm ON wm.workspace_id = w.id").
		Joins("LEFT JOIN (?) m ON m.workspace_id = w.id", memberSub).
		Joins("LEFT JOIN (?) p ON p.workspace_id = w.id", projectSub).
		Where("wm.user_id = ?", userID).
		Group("w.id, m.member_count, p.project_count").
		Scan(&result).Error

	if err != nil {
		logger.Logger.Error("error getting all workspaces", zap.Error(err))
		return nil, err
	}

	return &workspaceDto.WorkspaceList{Workspaces: result}, nil
}

func (s *service) CreateWorkspace(c context.Context, userID uuid.UUID, name string) (*workspaceDto.WorkspaceWithCounts, error) {
	var workspace = &models.Workspace{
		Name: name,
		Slug: utils.GenerateUniqueSlug(s.DB, name),
	}

	err := s.DB.Create(workspace).Error
	if err != nil {
		logger.Logger.Error("error creating workspace", zap.Error(err))
		return nil, err
	}

	var workspaceMember = &models.WorkspaceMember{
		WorkspaceID: workspace.ID,
		UserID:      userID,
		Role:        models.WorkspaceRoleOwner,
	}
	err = s.DB.Create(workspaceMember).Error
	if err != nil {
		logger.Logger.Error("error creating workspace_member", zap.Error(err))
		return nil, err
	}

	return &workspaceDto.WorkspaceWithCounts{
		WorkspaceResponse: workspaceDto.WorkspaceResponse{
			ID:          workspace.ID,
			Name:        workspace.Name,
			Slug:        workspace.Slug,
			BillingPlan: workspace.BillingPlan,
		},
		MemberCount:  0,
		ProjectCount: 0,
	}, nil
}

func (s *service) GetWorkspace(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*workspaceDto.WorkspaceDetailResponse, error) {
	var workspace models.Workspace

	err := s.DB.
		Preload("Projects").
		Preload("Members.User").
		First(&workspace, "id = ?", workspaceID).Error
	if err != nil {
		logger.Logger.Error("error getting workspace", zap.Error(err))
		return nil, err
	}

	var projects []projectDto.ProjectDetails
	var workspaceMembers []workspaceMemberDto.WorkspaceMemberDetail

	for _, project := range workspace.Projects {
		projects = append(projects, projectDto.ProjectDetails{
			ID:              project.ID,
			RepositoryID:    project.RepositoryID,
			WorkspaceID:     project.WorkspaceID,
			Name:            project.Name,
			Framework:       project.Framework,
			BuildCommand:    project.BuildCommand,
			OutputDirectory: project.OutputDirectory,
			RootDirectory:   project.RootDirectory,
		})
	}
	for _, member := range workspace.Members {
		workspaceMembers = append(workspaceMembers, workspaceMemberDto.WorkspaceMemberDetail{
			WorkspaceUser: userDto.WorkspaceUser{
				ID:             member.User.ID,
				Name:           member.User.Name,
				Email:          member.User.Email,
				ProfilePicture: member.User.ProfilePicture,
			},
			Role: member.Role,
		})
	}
	var res = &workspaceDto.WorkspaceDetailResponse{
		ID:          workspace.ID,
		Name:        workspace.Name,
		Slug:        workspace.Slug,
		BillingPlan: workspace.BillingPlan,
		Members:     workspaceMembers,
		Projects:    projects,
	}

	return res, nil
}

func (s *service) checkUserRole(userID uuid.UUID, workspaceID uuid.UUID) (models.WorkspaceRole, error) {
	var member models.WorkspaceMember
	err := s.DB.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&member).Error
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

func (s *service) UpdateWorkspace(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, req *workspaceDto.UpdateWorkspaceRequest) (*workspaceDto.WorkspaceResponse, error) {
	role, err := s.checkUserRole(userID, workspaceID)
	if err != nil {
		return nil, errors.New("unauthorized or workspace not found")
	}

	if role != models.WorkspaceRoleOwner && role != models.WorkspaceRoleAdmin {
		return nil, errors.New("action not allowed: must be owner or admin")
	}

	var workspace models.Workspace
	err = s.DB.First(&workspace, "id = ?", workspaceID).Error
	if err != nil {
		return nil, err
	}

	workspace.Name = req.Name
	workspace.Slug = utils.GenerateUniqueSlug(s.DB, req.Name)

	if err := s.DB.Save(&workspace).Error; err != nil {
		return nil, err
	}

	return &workspaceDto.WorkspaceResponse{
		ID:          workspace.ID,
		Name:        workspace.Name,
		Slug:        workspace.Slug,
		BillingPlan: workspace.BillingPlan,
	}, nil
}

func (s *service) DeleteWorkspace(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) error {
	role, err := s.checkUserRole(userID, workspaceID)
	if err != nil {
		return errors.New("unauthorized or workspace not found")
	}

	if role != models.WorkspaceRoleOwner {
		return errors.New("action not allowed: must be owner")
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("workspace_id = ?", workspaceID).Delete(&models.WorkspaceMember{}).Error; err != nil {
			return err
		}
		if err := tx.Where("workspace_id = ?", workspaceID).Delete(&models.Project{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", workspaceID).Delete(&models.Workspace{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *service) ListWorkspaceMembers(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) ([]workspaceMemberDto.WorkspaceMemberDetail, error) {
	_, err := s.checkUserRole(userID, workspaceID)
	if err != nil {
		return nil, errors.New("unauthorized or workspace not found")
	}

	var members []models.WorkspaceMember
	if err := s.DB.Preload("User").Where("workspace_id = ?", workspaceID).Find(&members).Error; err != nil {
		return nil, err
	}

	var res []workspaceMemberDto.WorkspaceMemberDetail
	for _, member := range members {
		res = append(res, workspaceMemberDto.WorkspaceMemberDetail{
			WorkspaceUser: userDto.WorkspaceUser{
				ID:             member.User.ID,
				Name:           member.User.Name,
				Email:          member.User.Email,
				ProfilePicture: member.User.ProfilePicture,
			},
			Role: member.Role,
		})
	}
	return res, nil
}

func (s *service) UpdateWorkspaceMemberRole(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, targetUserID uuid.UUID, role models.WorkspaceRole) error {
	currentUserRole, err := s.checkUserRole(userID, workspaceID)
	if err != nil {
		return errors.New("unauthorized or workspace not found")
	}

	if currentUserRole != models.WorkspaceRoleOwner && currentUserRole != models.WorkspaceRoleAdmin {
		return errors.New("action not allowed: must be owner or admin")
	}

	targetUserRole, err := s.checkUserRole(targetUserID, workspaceID)
	if err != nil {
		return errors.New("target user not found in workspace")
	}

	if currentUserRole == models.WorkspaceRoleAdmin {
		if targetUserRole == models.WorkspaceRoleOwner || role == models.WorkspaceRoleOwner {
			return errors.New("action not allowed: admins cannot modify owners or promote to owner")
		}
	}

	if userID == targetUserID {
		return errors.New("action not allowed: cannot modify your own role")
	}

	return s.DB.Model(&models.WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, targetUserID).
		Update("role", role).Error
}

func (s *service) RemoveWorkspaceMember(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, targetUserID uuid.UUID) error {
	currentUserRole, err := s.checkUserRole(userID, workspaceID)
	if err != nil {
		return errors.New("unauthorized or workspace not found")
	}

	if currentUserRole != models.WorkspaceRoleOwner && currentUserRole != models.WorkspaceRoleAdmin {
		return errors.New("action not allowed: must be owner or admin")
	}

	targetUserRole, err := s.checkUserRole(targetUserID, workspaceID)
	if err != nil {
		return errors.New("target user not found in workspace")
	}

	if currentUserRole == models.WorkspaceRoleAdmin {
		if targetUserRole == models.WorkspaceRoleOwner {
			return errors.New("action not allowed: admins cannot remove owners")
		}
	}

	if userID == targetUserID {
		return errors.New("action not allowed: cannot remove yourself, please leave instead")
	}

	return s.DB.Where("workspace_id = ? AND user_id = ?", workspaceID, targetUserID).Delete(&models.WorkspaceMember{}).Error
}

func (s *service) ListIntegrations(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) ([]workspaceDto.IntegrationResponse, error) {
	var integrations []models.Integration

	if err := s.DB.Preload("Workspace").Where("workspace_id = ?", workspaceID).Find(&integrations).Error; err != nil {
		logger.Logger.Error("list integrations failed", zap.Error(err))
		return nil, err
	}

	var res []workspaceDto.IntegrationResponse
	for _, integration := range integrations {
		res = append(res, workspaceDto.IntegrationResponse{
			ID:          integration.ID,
			WorkspaceID: integration.WorkspaceID,
			Provider:    integration.Provider,
			ProviderID:  integration.ProviderID,
			ExpiresAt:   integration.ExpiresAt,
		})
	}

	return res, nil
}

func (s *service) ConnectGithub(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) (string, error) {
	state := fmt.Sprintf("%s:%s", utils.GenerateRandomString(16), workspaceID)
	url := config.GithubOauthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	return url, nil
}

func (s *service) GithubCallback(c context.Context, userID uuid.UUID, code string, state string) (*workspaceDto.IntegrationResponse, error) {
	parts := strings.Split(state, ":")
	if len(parts) < 2 {
		return nil, errors.New("invalid github state")
	}
	workspaceID, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, errors.New("invalid github workspace")
	}

	token, _ := config.GithubOauthConfig.Exchange(c, code)

	logger.Logger.Info("token", zap.Any("token", token))

	integration := models.Integration{
		WorkspaceID:  workspaceID,
		Provider:     "github",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}
	err = s.DB.Create(&integration).Error
	if err != nil {
		return nil, err
	}

	res := &workspaceDto.IntegrationResponse{
		ID:          integration.ID,
		WorkspaceID: integration.WorkspaceID,
		Provider:    integration.Provider,
		ProviderID:  integration.ProviderID,
		ExpiresAt:   integration.ExpiresAt,
	}

	return res, nil
}

func (s *service) DeleteIntegration(c context.Context, userID uuid.UUID, integrationID uuid.UUID) error {
	return nil
}

func (s *service) ListWorkspaceRepos(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) ([]workspaceDto.GithubRepoResponse, error) {
	var integration models.Integration

	err := s.DB.Where("workspace_id = ? AND provider = ?", workspaceID, models.IntegrationProviderGithub).
		Order("created_at DESC").
		First(&integration).
		Error
	if err != nil {
		logger.Logger.Error("fetch integration", zap.Error(err))
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: integration.AccessToken},
	)
	tc := oauth2.NewClient(c, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryListByAuthenticatedUserOptions{
		Affiliation: "owner,collaborator,organization_member",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	repos, _, err := client.Repositories.ListByAuthenticatedUser(c, opts)
	if err != nil {
		logger.Logger.Error("list repos failed", zap.Error(err))
		return nil, errors.New("list repos failed")
	}

	var response []workspaceDto.GithubRepoResponse
	for _, repo := range repos {
		response = append(response, workspaceDto.GithubRepoResponse{
			ID:            repo.GetID(),
			Name:          repo.GetName(),
			FullName:      repo.GetFullName(),
			Private:       repo.GetPrivate(),
			URL:           repo.GetURL(),
			Description:   repo.GetDescription(),
			UpdatedAt:     repo.GetUpdatedAt(),
			DefaultBranch: repo.GetDefaultBranch(),
		})
	}

	return response, nil
}

func (s *service) CreateProject(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, repositoryID int64, name string, repoFullName string, desc string, url string, private bool, buildCommand string, outputDir string, deployBranch string, rootDir string) (*workspaceDto.ProjectData, error) {
	var repo = &models.Repository{
		WorkspaceID:   workspaceID,
		Provider:      models.IntegrationProviderGithub,
		ExternalID:    strconv.FormatInt(repositoryID, 10),
		Name:          name,
		FullName:      repoFullName,
		Description:   desc,
		URL:           url,
		DefaultBranch: deployBranch,
		IsPrivate:     private,
	}
	err := s.DB.Create(&repo).Error
	if err != nil {
		logger.Logger.Error("create repo failed", zap.Error(err))
		return nil, err
	}

	var project = &models.Project{
		WorkspaceID:     workspaceID,
		RepositoryID:    repo.ID,
		Name:            name,
		RepoFullName:    repoFullName,
		Framework:       "",
		BuildCommand:    buildCommand,
		OutputDirectory: outputDir,
		RootDirectory:   rootDir,
		DeployBranch:    deployBranch,
		WebhookSecret:   utils.GenerateRandomString(32),
	}
	err = s.DB.Create(&project).Error
	if err != nil {
		logger.Logger.Error("create project failed", zap.Error(err))
		return nil, err
	}

	ghconfig := &github.HookConfig{
		URL:         github.String("https://af1e-106-211-50-211.ngrok-free.app/webhooks/github"),
		ContentType: github.String("json"),
		Secret:      github.String(project.WebhookSecret), // Decrypted secret from your DB
	}

	hook := &github.Hook{
		Name:   github.String("web"),
		Config: ghconfig,
		Events: []string{"push"},
		Active: github.Bool(true),
	}

	var integration models.Integration

	err = s.DB.Where("workspace_id = ? AND provider = ?", workspaceID, models.IntegrationProviderGithub).
		Order("created_at DESC").
		First(&integration).
		Error
	if err != nil {
		logger.Logger.Error("fetch integration", zap.Error(err))
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: integration.AccessToken},
	)
	tc := oauth2.NewClient(c, ts)
	client := github.NewClient(tc)

	parts := strings.Split(project.RepoFullName, "/")
	owner := parts[0]
	repoName := parts[1]

	_, _, err = client.Repositories.CreateHook(c, owner, repoName, hook)
	if err != nil {
		logger.Logger.Error("create repo hook failed", zap.Error(err))
		return nil, err
	}

	var res = &workspaceDto.ProjectData{
		ID:              project.ID,
		WorkspaceID:     project.WorkspaceID,
		RepositoryID:    project.RepositoryID,
		Name:            project.Name,
		RepoFullName:    project.RepoFullName,
		Framework:       project.Framework,
		BuildCommand:    project.BuildCommand,
		OutputDirectory: project.OutputDirectory,
		DeployBranch:    project.DeployBranch,
		RootDirectory:   project.RootDirectory,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
	}
	return res, nil
}
