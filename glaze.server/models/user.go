package models

type User struct {
	Base
	Name           string
	Email          string `gorm:"uniqueIndex"`
	Password       string
	ProfilePicture string
	Role           string `gorm:"default:user"`

	Workspaces []WorkspaceMember `gorm:"foreignKey:UserID"`
}
