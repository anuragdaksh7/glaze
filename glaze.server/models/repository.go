package models

import (
	"github.com/google/uuid"
)

type Repository struct {
	Base
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null"`
	Provider    string    // "github" or "gitlab"
	ExternalID  string    `gorm:"not null"` // The ID from the git provider
	FullName    string    `gorm:"not null"` // e.g., "owner/repo"
	Project     Project
}
