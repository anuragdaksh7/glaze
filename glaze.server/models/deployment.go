package models

import (
	"time"

	"github.com/google/uuid"
)

type DeploymentStatus string

const (
	StatusQueued    DeploymentStatus = "queued"
	StatusCloning   DeploymentStatus = "cloning"
	StatusBuilding  DeploymentStatus = "building"
	StatusSuccess   DeploymentStatus = "success"
	StatusFailed    DeploymentStatus = "failed"
	StatusCancelled DeploymentStatus = "cancelled"
)

type Deployment struct {
	Base
	ProjectID uuid.UUID `json:"projectId" gorm:"type:uuid;not null;index"`
	Project   Project   `json:"-" gorm:"foreignKey:ProjectID"`

	// Git Metadata
	CommitHash string `json:"commitHash" gorm:"not null"`
	CommitMsg  string `json:"commitMsg"`
	AuthorName string `json:"authorName"`
	Branch     string `json:"branch" gorm:"not null"`

	ContainerID  string `gorm:"type:text" json:"container_id"` // <--- ADD THIS
	ExternalPort int    `json:"external_port"`                 // Useful for mapping localhost:XXXX
	ImageName    string `json:"image_name"`

	// Lifecycle Status
	// Status options: queued, cloning, building, success, failed, cancelled
	Status DeploymentStatus `json:"status" gorm:"type:deployment_status;default:'queued';index"`

	// Logs & Debugging
	// We use 'text' type for logs because they can get very long
	Logs string `json:"logs" gorm:"type:text"`

	// Metrics
	BuildDuration int64      `json:"buildDuration"` // in seconds
	StartedAt     *time.Time `json:"startedAt"`
	FinishedAt    *time.Time `json:"finishedAt"`
}
