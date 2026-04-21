package models

import (
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	WorkspaceID uint   `gorm:"not null"`
	Provider    string // "github" or "gitlab"
	ExternalID  string `gorm:"not null"` // The ID from the git provider
	FullName    string `gorm:"not null"` // e.g., "owner/repo"
	Project     Project
}
