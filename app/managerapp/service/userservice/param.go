package userservice

import (
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/types"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	User  User                               `json:"user"`
	Token tokenservice.TokenGenerateResponse `json:"token"`
}

type UserCreateRequest struct {
	ID       types.ID     `json:"-"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	Fullname string       `json:"fullname"`
	Roles    []types.Role `json:"roles"`
}

type UserCreateResponse struct {
	User
}

type UserProfileResponse struct {
	User
}
