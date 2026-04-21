package config

import "track/models"

func SyncDB() {
	err := DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return
	}
}
