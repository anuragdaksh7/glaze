package workspaceMemberDto

import (
	userDto "glaze/dto/user"
	"glaze/models"
)

type WorkspaceMemberDetail struct {
	userDto.WorkspaceUser
	Role models.WorkspaceRole `json:"role"`
}
