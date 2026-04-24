package config

import "glaze/models"

func SyncDB() {
	DB.Exec(`
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'workspace_role'
    ) THEN
        CREATE TYPE workspace_role AS ENUM (
            'owner','admin','member','viewer'
        );
    END IF;
END $$;
`)
	DB.Exec(`
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'integration_provider'
    ) THEN
        CREATE TYPE integration_provider AS ENUM (
            'github'
        );
    END IF;
END $$;
`)
	err := DB.AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.Project{},
		&models.WorkspaceMember{},
		&models.Repository{},
		&models.Deployment{},
		&models.Domain{},
		&models.EnvVar{},
		&models.Integration{},
	)
	if err != nil {
		return
	}
}
