package models

import (
	"gorm.io/gorm"
)

type EnvVar struct {
	gorm.Model
	ProjectID uint   `gorm:"not null"`
	Key       string `gorm:"not null"`
	Value     string `gorm:"not null"` // Should be stored encrypted
	Scope     string `gorm:"default:'all'"` // "production", "preview", or "all"
}
