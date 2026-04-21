package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"glaze/config"
	userDto "glaze/dto/user"
	"glaze/logger"
	"glaze/models"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	timeout time.Duration
	DB      *gorm.DB
}

func NewService() Service {
	return &service{
		time.Duration(20) * time.Second,
		config.DB,
	}
}

func (s *service) SignUp(c context.Context, req *userDto.SignUpReq) (*userDto.SignUpResp, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		logger.Logger.Error("SignUpReq fields: ", zap.Any("req", req))
		return nil, errors.New("invalid fields")
	}
	var user models.User
	s.DB.Where("email = ?", req.Email).First(&user)
	if user.ID != 0 {
		logger.Logger.Error("User already exists")
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error hashing password: %s", err))
		return nil, err
	}

	user = models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hash),
	}
	res := s.DB.Create(&user)
	if res.Error != nil {
		logger.Logger.Error("Error creating user: ", zap.Error(res.Error))
		return nil, res.Error
	}
	logger.Logger.Info("User created: ", zap.String("email", user.Email))

	resp := &userDto.SignUpResp{
		Id: user.ID,
	}

	return resp, nil
}

func (s *service) Login(c context.Context, req *userDto.LoginReq) (*userDto.LoginResp, error) {
	if req.Email == "" || req.Password == "" {
		logger.Logger.Error("SignUpReq fields: ", zap.Any("req", req))
		return nil, errors.New("invalid fields")
	}
	_config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	s.DB.First(&user, "email = ?", req.Email)

	if user.ID == 0 {
		logger.Logger.Error("User not found: ", zap.String("email", req.Email))
		return nil, errors.New("user not found")
	}

	fmt.Println(req.Password, _config.MasterPassword)
	if req.Password == _config.MasterPassword {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		tokenString, err := token.SignedString([]byte(_config.JwtSecret))
		if err != nil {
			logger.Logger.Error("Error creating token: ", zap.String("token", tokenString))
			return nil, err
		}

		return &userDto.LoginResp{
			tokenString,
			user.ID,
			user.Name,
			user.Email,
			user.Role,
			user.ProfilePicture,
		}, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logger.Logger.Error("User found: ", zap.String("email", req.Email))
		return nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(_config.JwtSecret))
	if err != nil {
		logger.Logger.Error("Error creating token: ", zap.String("token", tokenString))
		return nil, err
	}

	return &userDto.LoginResp{
		tokenString,
		user.ID,
		user.Name,
		user.Email,
		user.Role,
		user.ProfilePicture,
	}, nil
}

func (s *service) Me(c context.Context, userId uint) (*userDto.GetMeRes, error) {
	//createDefaultLinks(s)
	var user models.User

	s.DB.First(&user, "id = ?", userId)
	if user.ID == 0 {
		logger.Logger.Error("User not found: ", zap.String("email", user.Email))
		return nil, errors.New("user not found")
	}

	res := &userDto.GetMeRes{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Role:           user.Role,
	}

	return res, nil
}
