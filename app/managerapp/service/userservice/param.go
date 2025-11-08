package userservice

import (
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
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
	ID          types.ID     `json:"-"`
	Username    string       `json:"username"`
	Password    string       `json:"password"`
	Fullname    string       `json:"fullname"`
	Email       string       `json:"email"`
	PhoneNumber string       `json:"phone_number"`
	Roles       []types.Role `json:"roles"`
}

type UserCreateResponse struct {
	User
}

type UserUpdateFromSuperuserRequest struct {
	Username    string       `json:"username"`
	Fullname    string       `json:"fullname"`
	Email       string       `json:"email"`
	PhoneNumber string       `json:"phone_number"`
	Roles       []types.Role `json:"roles"`
}

type UserUpdateFromOwnRequest struct {
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type UserUpdateResponse struct {
	User
}

type UserProfileResponse struct {
	User
}

type ListUserRequest struct {
	Username  string               `json:"username"`
	Roles     []types.Role         `json:"roles"`
	Paginated paginate.RequestBase `json:"paginated"`
}

type ListUserResponse struct {
	Results  []User                `json:"results"`
	Paginate paginate.ResponseBase `json:"paginate"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type RegisterGuestUserRequest struct {
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type RegisterGuestUserResponse struct {
	User

	QCToken string `json:"qc_token"`
}

type UpdateGuestUserRequest struct {
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateGuestUserResponse struct {
	User
}

type IdentifyClientRequest struct {
	ExternalUserID string `json:"user_id"`
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
}

type IdentifyClientResponse struct {
	User

	QCToken string `json:"qc_token"`
}
