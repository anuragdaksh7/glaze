package config

import "glaze/models"

func SyncDB() {
	err := DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return
	}
}
