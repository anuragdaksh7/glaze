package workspaceDto

import (
	projectDto "glaze/dto/project"
	workspaceMemberDto "glaze/dto/workspacemember"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/google/uuid"
)

type WorkspaceResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	BillingPlan string    `json:"billing_plan"`
}

type WorkspaceDetailResponse struct {
	ID          uuid.UUID                                  `json:"id"`
	Name        string                                     `json:"name"`
	Slug        string                                     `json:"slug"`
	BillingPlan string                                     `json:"billing_plan"`
	Members     []workspaceMemberDto.WorkspaceMemberDetail `json:"members"`
	Projects    []projectDto.ProjectDetails                `json:"projects"`
}

type WorkspaceWithCounts struct {
	WorkspaceResponse
	MemberCount  int `json:"member_count"`
	ProjectCount int `json:"project_count"`
}

type WorkspaceList struct {
	Workspaces []WorkspaceWithCounts `json:"workspaces"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

type UpdateWorkspaceRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateWorkspaceMemberRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type WorkspaceMemberReq struct {
	WorkspaceID string `uri:"workspace_id" binding:"required,uuid"`
	UserID      string `uri:"user_id" binding:"required,uuid"`
}

type GetWorkspaceByIDReq struct {
	ID string `uri:"workspace_id" binding:"required,uuid"`
}

type IntegrationResponse struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Provider    string    `json:"provider"`
	ProviderID  string    `json:"provider_id"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type GithubRepoResponse struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`
	FullName      string           `json:"fullName"`
	Private       bool             `json:"private"`
	URL           string           `json:"url"`
	Description   string           `json:"description"`
	UpdatedAt     github.Timestamp `json:"updatedAt"`
	DefaultBranch string           `json:"defaultBranch"`
}

type CreateProjectRequest struct {
	RepositoryID int64  `json:"repositoryId" binding:"required"`
	RepoFullName string `json:"repoFullName" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	URL          string `json:"url"`
	IsPrivate    bool   `json:"isPrivate"`

	BuildCommand    string `json:"buildCommand"`    // Default: 'npm run build'
	OutputDirectory string `json:"outputDirectory"` // Default: 'dist'
	DeployBranch    string `json:"deployBranch"`    // Default: 'main'
	RootDirectory   string `json:"rootDirectory"`   // Default: '/'
}

type ProjectData struct {
	ID              uuid.UUID `json:"id"`
	WorkspaceID     uuid.UUID `json:"workspaceId"`
	RepositoryID    uuid.UUID `json:"repositoryId"`
	Name            string    `json:"name"`
	RepoFullName    string    `json:"repoFullName"`
	Framework       string    `json:"framework"`
	BuildCommand    string    `json:"buildCommand"`
	OutputDirectory string    `json:"outputDirectory"`
	DeployBranch    string    `json:"deployBranch"`
	RootDirectory   string    `json:"rootDirectory"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
