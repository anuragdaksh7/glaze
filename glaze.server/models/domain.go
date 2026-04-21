package models

import (
	"gorm.io/gorm"
)

type Domain struct {
	gorm.Model
	ProjectID uint   `gorm:"not null"`
	Hostname  string `gorm:"uniqueIndex;not null"`
	IsCustom  bool   `gorm:"default:false"`
	SSLStatus string `gorm:"default:'pending'"`
}
