package adminDto

type GetAllUsersRes struct {
	Id          uint   `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Public      int64  `json:"public"`
	Collections int64  `json:"collections"`
	Bookmarks   int64  `json:"bookmarks"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
}

type GetAdminAccessReq struct {
	UserId uint `json:"user_id"`
}
