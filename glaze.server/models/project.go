package models

import (
	"github.com/google/uuid"
)

type Project struct {
	Base
	WorkspaceID  uuid.UUID `json:"workspaceId" gorm:"type:uuid;not null;index"`
	RepositoryID uuid.UUID `json:"repositoryId" gorm:"type:uuid;not null"`

	// Descriptive Info
	Name         string `json:"name" gorm:"not null"`
	RepoFullName string `json:"repoFullName" gorm:"not null"` // e.g., "anurag/my-cool-app"
	Framework    string `json:"framework"`                    // nextjs, vite, etc.

	// Build Settings
	BuildCommand    string `json:"buildCommand" gorm:"default:'npm run build'"`
	OutputDirectory string `json:"outputDirectory" gorm:"default:'dist'"`
	RootDirectory   string `json:"rootDirectory" gorm:"default:'/'"`
	DeployBranch    string `json:"deployBranch" gorm:"default:'main'"`

	// Security & Automation
	// This secret is used to verify GitHub Webhook signatures
	WebhookSecret string `json:"-" gorm:"not null"`

	// Relationships
	Deployments []Deployment `json:"deployments,omitempty"`
	EnvVars     []EnvVar     `json:"envVars,omitempty"`
	Domains     []Domain     `json:"domains,omitempty"`
}
