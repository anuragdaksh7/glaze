package models

import (
	"github.com/google/uuid"
)

type Repository struct {
	Base
	WorkspaceID uuid.UUID `json:"workspaceId" gorm:"type:uuid;not null;index"`

	// Identity from Provider
	Provider   IntegrationProvider `json:"provider" gorm:"type:integration_provider;default:'github'"`
	ExternalID string              `json:"externalId" gorm:"uniqueIndex;not null"` // GitHub's unique ID

	// Meta Info (Cached from GitHub)
	Name          string `json:"name" gorm:"not null"`
	FullName      string `json:"fullName" gorm:"not null"` // owner/repo
	Description   string `json:"description"`
	URL           string `json:"url"`           // https://github.com/owner/repo
	DefaultBranch string `json:"defaultBranch"` // used to suggest the 'DeployBranch'
	IsPrivate     bool   `json:"isPrivate"`

	ProjectID *uuid.UUID `json:"projectId" gorm:"type:uuid"`
	//Project   *Project   `json:"project,omitempty" gorm:"constraint:OnDelete:SET NULL;"`
	//Workspace Workspace  `json:"workspace" gorm:"constraint:OnDelete:SET NULL;"`
}
