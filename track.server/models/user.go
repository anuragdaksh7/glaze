package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name             string
	Email            string `gorm:"uniqueIndex"`
	Password         string
	ProfilePicture   string
	DefaultProfileID *uint
	Role             string `gorm:"default:user"`
	otp              string
	//Profiles          []Profile            `gorm:"foreignKey:UserID;references:ID"`
	//Collections       []Collection         `gorm:"foreignKey:UserID;references:ID"`
	//SharedCollections []CollectionShare    `gorm:"foreignKey:UserID;references:ID"`
	//Bookmarks         []Bookmark           `gorm:"foreignKey:UserID;references:ID"`
	//UserLabels        []UserLabel          `gorm:"foreignKey:UserID;references:ID"`
	//BookmarkLabels    []BookmarkLabel      `gorm:"foreignKey:UserID;references:ID"`
	//Statuses          []UserBookmarkStatus `gorm:"foreignKey:UserID;references:ID"`
	//Invites           []Invite             `gorm:"foreignKey:UserID;references:ID"`
	//CreditLogs        []CreditLog          `gorm:"foreignKey:UserID;references:ID"`
	//UserCredits       *UserCredits         `gorm:"foreignKey:UserID;references:ID"`
	//Payments          []Payment            `gorm:"foreignKey:UserID;references:ID"`
}
