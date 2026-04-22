package projectDto

import "github.com/google/uuid"

type ProjectDetails struct {
	ID              uuid.UUID `json:"id"`
	RepositoryID    uuid.UUID `json:"repository_id"`
	WorkspaceID     uuid.UUID `json:"workspace_id"`
	Name            string    `json:"name"`
	Framework       string    `json:"framework"`
	BuildCommand    string    `json:"build_command"`
	OutputDirectory string    `json:"output_dir"`
	RootDirectory   string    `json:"root_dir"`
}
