package userDto

type SignUpReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResp struct {
	Id uint `json:"id"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token          string `json:"token"`
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	ProfilePicture string `json:"profilePicture"`
}

type GetMeRes struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	Role           string `json:"role"`
}
