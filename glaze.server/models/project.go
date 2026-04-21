package models

import (
	"github.com/google/uuid"
)

type Project struct {
	Base
	WorkspaceID     uuid.UUID `gorm:"type:uuid;not null"`
	RepositoryID    uuid.UUID `gorm:"type:uuid;not null"`
	Name            string    `gorm:"not null"`
	Framework       string    // e.g., "nextjs", "vite"
	BuildCommand    string    `gorm:"default:'npm run build'"`
	OutputDirectory string    `gorm:"default:'dist'"`
	RootDirectory   string    `gorm:"default:'/'"`
	Deployments     []Deployment
	EnvVars         []EnvVar
	Domains         []Domain
}
