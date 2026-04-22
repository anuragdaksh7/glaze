package workspace

import (
	"context"
	workspaceDto "glaze/dto/workspace"

	"github.com/google/uuid"
)

type Service interface {
	GetAllWorkspaces(c context.Context, userID uuid.UUID) (*workspaceDto.WorkspaceList, error)
	CreateWorkspace(c context.Context, userID uuid.UUID, name string) (*workspaceDto.WorkspaceWithCounts, error)
	GetWorkspace(c context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*workspaceDto.WorkspaceDetailResponse, error)
}
