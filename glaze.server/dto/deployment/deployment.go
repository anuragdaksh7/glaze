package deploymentDto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDeploymentRequest struct {
	Branch    string `json:"branch"`
	CommitSha string `json:"commitSha"`
}

type DeploymentResponse struct {
	ID            uuid.UUID  `json:"id"`
	ProjectID     uuid.UUID  `json:"projectId"`
	Status        string     `json:"status"`
	URL           *string    `json:"url"`
	CommitSha     *string    `json:"commitSha"`
	CommitMessage *string    `json:"commitMessage"`
	AuthorName    *string    `json:"authorName"`
	Branch        string     `json:"branch"`
	Logs          string     `json:"logs"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	FinishedAt    *time.Time `json:"finishedAt"`
	DurationMs    *int64     `json:"durationMs"`
}

type DeploymentByIDReq struct {
	WorkspaceID  string `uri:"workspace_id" binding:"required,uuid"`
	ProjectID    string `uri:"project_id" binding:"required,uuid"`
	DeploymentID string `uri:"deployment_id" binding:"required,uuid"`
}

type DeploymentsByProjectReq struct {
	WorkspaceID string `uri:"workspace_id" binding:"required,uuid"`
	ProjectID   string `uri:"project_id" binding:"required,uuid"`
}
