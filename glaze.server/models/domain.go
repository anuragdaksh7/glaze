package models

import (
	"github.com/google/uuid"
)

type Domain struct {
	Base
	ProjectID uuid.UUID `gorm:"type:uuid;not null"`
	Hostname  string    `gorm:"uniqueIndex;not null"`
	IsCustom  bool      `gorm:"default:false"`
	SSLStatus string    `gorm:"default:'pending'"`
}
