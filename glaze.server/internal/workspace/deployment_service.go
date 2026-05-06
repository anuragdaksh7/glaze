package workspace

import (
	"context"
	"errors"
	"fmt"
	deploymentDto "glaze/dto/deployment"
	"glaze/internal/tasks"
	"glaze/logger"
	"glaze/models"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func deploymentURL(repoFullName string) *string {
	slug := strings.ToLower(repoFullName)
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, ".", "-")
	url := fmt.Sprintf("https://%s.localhost", slug)
	return &url
}

func mapDeploymentStatus(status models.DeploymentStatus) string {
	switch status {
	case models.StatusSuccess:
		return "ready"
	case models.StatusFailed:
		return "error"
	case models.StatusCancelled:
		return "canceled"
	case models.StatusCloning, models.StatusBuilding:
		return "building"
	default:
		return "queued"
	}
}

func deploymentDurationMs(d *models.Deployment) *int64 {
	if d.BuildDuration > 0 {
		ms := d.BuildDuration * 1000
		return &ms
	}
	if d.StartedAt != nil && d.FinishedAt != nil {
		ms := d.FinishedAt.Sub(*d.StartedAt).Milliseconds()
		return &ms
	}
	return nil
}

func deploymentResponse(project *models.Project, d *models.Deployment) *deploymentDto.DeploymentResponse {
	commitSha := d.CommitHash
	commitMsg := d.CommitMsg
	authorName := d.AuthorName

	var commitShaPtr *string
	if commitSha != "" {
		commitShaPtr = &commitSha
	}

	var commitMsgPtr *string
	if commitMsg != "" {
		commitMsgPtr = &commitMsg
	}

	var authorNamePtr *string
	if authorName != "" {
		authorNamePtr = &authorName
	}

	return &deploymentDto.DeploymentResponse{
		ID:            d.ID,
		ProjectID:     d.ProjectID,
		Status:        mapDeploymentStatus(d.Status),
		URL:           deploymentURL(project.RepoFullName),
		CommitSha:     commitShaPtr,
		CommitMessage: commitMsgPtr,
		AuthorName:    authorNamePtr,
		Branch:        d.Branch,
		Logs:          d.Logs,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
		FinishedAt:    d.FinishedAt,
		DurationMs:    deploymentDurationMs(d),
	}
}

func (s *service) workspaceProject(userID, workspaceID, projectID uuid.UUID) (*models.Project, error) {
	if _, err := s.checkUserRole(userID, workspaceID); err != nil {
		return nil, errors.New("unauthorized or workspace not found")
	}

	var project models.Project
	if err := s.DB.Where("id = ? AND workspace_id = ?", projectID, workspaceID).First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *service) githubClientForWorkspace(c context.Context, workspaceID uuid.UUID) (*github.Client, error) {
	var integration models.Integration
	if err := s.DB.Where("workspace_id = ? AND provider = ?", workspaceID, models.IntegrationProviderGithub).
		Order("created_at DESC").
		First(&integration).Error; err != nil {
		logger.Logger.Error("fetch integration", zap.Error(err))
		return nil, err
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: integration.AccessToken})
	tc := oauth2.NewClient(c, ts)
	return github.NewClient(tc), nil
}

func (s *service) latestCommitForBranch(c context.Context, workspaceID uuid.UUID, repoFullName string, branch string, commitSha string) (string, string, string, error) {
	client, err := s.githubClientForWorkspace(c, workspaceID)
	if err != nil {
		return "", "", "", err
	}

	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		return "", "", "", errors.New("invalid repository full name")
	}

	owner := parts[0]
	repo := parts[1]

	if commitSha != "" {
		commit, _, err := client.Repositories.GetCommit(c, owner, repo, commitSha, nil)
		if err != nil {
			return "", "", "", err
		}
		return commit.GetSHA(), commit.GetCommit().GetMessage(), commit.GetCommit().GetAuthor().GetName(), nil
	}

	opts := &github.CommitsListOptions{
		SHA: branch,
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}

	commits, _, err := client.Repositories.ListCommits(c, owner, repo, opts)
	if err != nil {
		return "", "", "", err
	}
	if len(commits) == 0 {
		return "", "", "", errors.New("no commits found for branch")
	}

	commit := commits[0]
	return commit.GetSHA(), commit.GetCommit().GetMessage(), commit.GetCommit().GetAuthor().GetName(), nil
}

func (s *service) ListDeployments(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, projectID uuid.UUID) ([]deploymentDto.DeploymentResponse, error) {
	project, err := s.workspaceProject(userID, workspaceID, projectID)
	if err != nil {
		return nil, err
	}

	var deployments []models.Deployment
	if err := s.DB.Where("project_id = ?", project.ID).Order("created_at DESC").Find(&deployments).Error; err != nil {
		return nil, err
	}

	res := make([]deploymentDto.DeploymentResponse, 0, len(deployments))
	for i := range deployments {
		res = append(res, *deploymentResponse(project, &deployments[i]))
	}

	return res, nil
}

func (s *service) GetDeployment(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, projectID uuid.UUID, deploymentID uuid.UUID) (*deploymentDto.DeploymentResponse, error) {
	project, err := s.workspaceProject(userID, workspaceID, projectID)
	if err != nil {
		return nil, err
	}

	var deployment models.Deployment
	if err := s.DB.Where("id = ? AND project_id = ?", deploymentID, project.ID).First(&deployment).Error; err != nil {
		return nil, err
	}

	return deploymentResponse(project, &deployment), nil
}

func (s *service) TriggerDeployment(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, projectID uuid.UUID, req *deploymentDto.CreateDeploymentRequest) (*deploymentDto.DeploymentResponse, error) {
	project, err := s.workspaceProject(userID, workspaceID, projectID)
	if err != nil {
		return nil, err
	}

	branch := strings.TrimSpace(req.Branch)
	if branch == "" {
		branch = project.DeployBranch
	}

	commitSha, commitMsg, authorName, err := s.latestCommitForBranch(c, workspaceID, project.RepoFullName, branch, strings.TrimSpace(req.CommitSha))
	if err != nil {
		return nil, err
	}

	deployment := &models.Deployment{
		ProjectID:  project.ID,
		CommitHash: commitSha,
		CommitMsg:  commitMsg,
		AuthorName: authorName,
		Branch:     branch,
		Status:     models.StatusQueued,
	}
	if err := s.DB.Create(deployment).Error; err != nil {
		return nil, err
	}

	task, err := tasks.NewBuildTask(deployment.ID.String(), project.RepoFullName, branch, commitSha)
	if err != nil {
		return nil, err
	}

	if _, err := s.AsynqClient.Enqueue(task); err != nil {
		logger.Logger.Error("enqueue deployment task failed", zap.Error(err))
		_ = s.DB.Model(&models.Deployment{}).
			Where("id = ?", deployment.ID).
			Updates(map[string]interface{}{
				"status":         models.StatusFailed,
				"logs":           err.Error(),
				"finished_at":    time.Now(),
				"build_duration": 0,
			}).Error
		return nil, err
	}

	return deploymentResponse(project, deployment), nil
}
