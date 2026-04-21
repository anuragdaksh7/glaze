package models

import (
	"gorm.io/gorm"
)

type Deployment struct {
	gorm.Model
	ProjectID     uint   `gorm:"not null;index"`
	CreatorID     uint   `gorm:"not null"`
	CommitSHA     string `gorm:"not null"`
	CommitMessage string
	Branch        string `gorm:"not null"`
	Environment   string `gorm:"default:'preview'"` // "production" or "preview"
	Status        string `gorm:"default:'QUEUED'"`  // QUEUED, BUILDING, READY, ERROR
	BuildLogURL   string
}
