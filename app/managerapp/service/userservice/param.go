package userservice

import "github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	User  User                               `json:"user"`
	Token tokenservice.TokenGenerateResponse `json:"token"`
}
