package workspaceDto

import (
	projectDto "glaze/dto/project"
	workspaceMemberDto "glaze/dto/workspacemember"

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

type GetWorkspaceByIDReq struct {
	ID string `uri:"workspace_id" binding:"required,uuid"`
}
