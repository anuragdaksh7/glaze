package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	WorkspaceID     uint   `gorm:"not null"`
	RepositoryID    uint   `gorm:"not null"`
	Name            string `gorm:"not null"`
	Framework       string // e.g., "nextjs", "vite"
	BuildCommand    string `gorm:"default:'npm run build'"`
	OutputDirectory string `gorm:"default:'dist'"`
	RootDirectory   string `gorm:"default:'/'"`
	Deployments     []Deployment
	EnvVars         []EnvVar
	Domains         []Domain
}
