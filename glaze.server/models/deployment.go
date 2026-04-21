package models

import (
	"github.com/google/uuid"
)

type Deployment struct {
	Base
	ProjectID     uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatorID     uuid.UUID `gorm:"type:uuid;not null"`
	CommitSHA     string    `gorm:"not null"`
	CommitMessage string
	Branch        string `gorm:"not null"`
	Environment   string `gorm:"default:'preview'"` // "production" or "preview"
	Status        string `gorm:"default:'QUEUED'"`  // QUEUED, BUILDING, READY, ERROR
	BuildLogURL   string
}
