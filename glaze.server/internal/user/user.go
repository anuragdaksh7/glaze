package user

import (
	"context"
	userDto "glaze/dto/user"
)

type Service interface {
	SignUp(c context.Context, req *userDto.SignUpReq) (*userDto.SignUpResp, error)
	Login(c context.Context, req *userDto.LoginReq) (*userDto.LoginResp, error)
	Me(c context.Context, userId uint) (*userDto.GetMeRes, error)
}
