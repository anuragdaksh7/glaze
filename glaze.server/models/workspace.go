package models

type Workspace struct {
	Base
	Name        string `gorm:"not null"`
	Slug        string `gorm:"uniqueIndex;not null"`
	BillingPlan string `gorm:"default:'free'"`
	Projects    []Project
	Members     []User `gorm:"many2many:workspace_members;"`
}
