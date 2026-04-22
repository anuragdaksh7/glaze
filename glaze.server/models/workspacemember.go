package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceRole string

const (
	WorkspaceRoleOwner  WorkspaceRole = "owner"
	WorkspaceRoleAdmin  WorkspaceRole = "admin"
	WorkspaceRoleMember WorkspaceRole = "member"
	WorkspaceRoleViewer WorkspaceRole = "viewer"
)

type WorkspaceMember struct {
	WorkspaceID uuid.UUID `gorm:"primaryKey"`
	UserID      uuid.UUID `gorm:"primaryKey"`

	Role      WorkspaceRole `gorm:"type:workspace_role;default:'member'"`
	Workspace Workspace     `gorm:"foreignKey:WorkspaceID"`
	User      User          `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
}
