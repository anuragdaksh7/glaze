package models

import (
	"github.com/google/uuid"
)

type EnvVar struct {
	Base
	ProjectID uuid.UUID `gorm:"type:uuid;not null"`
	Key       string    `gorm:"not null"`
	Value     string    `gorm:"not null"`      // Should be stored encrypted
	Scope     string    `gorm:"default:'all'"` // "production", "preview", or "all"
}
