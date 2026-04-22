package workspace

import (
	"context"
	"glaze/config"
	projectDto "glaze/dto/project"
	userDto "glaze/dto/user"
	workspaceDto "glaze/dto/workspace"
	workspaceMemberDto "glaze/dto/workspacemember"
	"glaze/logger"
	"glaze/models"
	"glaze/utils"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
