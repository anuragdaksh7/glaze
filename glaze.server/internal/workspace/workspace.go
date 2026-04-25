package workspace

import (
	"context"
	workspaceDto "glaze/dto/workspace"
	workspaceMemberDto "glaze/dto/workspacemember"

	"glaze/models"

	"github.com/google/uuid"
)

type Service interface {
	GetAllWorkspaces(c context.Context, userID uuid.UUID) (*workspaceDto.WorkspaceList, error)
	CreateWorkspace(c context.Context, userID uuid.UUID, name string) (*workspaceDto.WorkspaceWithCounts, error)
	GetWorkspace(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*workspaceDto.WorkspaceDetailResponse, error)
	UpdateWorkspace(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, req *workspaceDto.UpdateWorkspaceRequest) (*workspaceDto.WorkspaceResponse, error)
	DeleteWorkspace(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) error
	ListWorkspaceMembers(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) ([]workspaceMemberDto.WorkspaceMemberDetail, error)
	UpdateWorkspaceMemberRole(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, targetUserID uuid.UUID, role models.WorkspaceRole) error
	RemoveWorkspaceMember(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, targetUserID uuid.UUID) error
	ListIntegrations(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) ([]workspaceDto.IntegrationResponse, error)
	ConnectGithub(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) (string, error)
	GithubCallback(c context.Context, userID uuid.UUID, code string, state string) (*workspaceDto.IntegrationResponse, error)
	DeleteIntegration(c context.Context, userID uuid.UUID, integrationID uuid.UUID) error
	ListWorkspaceRepos(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) ([]workspaceDto.GithubRepoResponse, error)
	CreateProject(c context.Context, userID uuid.UUID, workspaceID uuid.UUID, repositoryID int64, name string, repoFullName string, desc string, url string, private bool, buildCommand string, outputDir string, deployBranch string, rootDir string) (*workspaceDto.ProjectData, error)
}
