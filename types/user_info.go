package types

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.RegisteredClaims

	UserID    ID        `json:"user_id"`
	Roles     []Role    `json:"roles"`
	TokenType TokenType `json:"token_type"`
}

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
	TokenTypeGuest   TokenType = "guest"
	TokenTypeClient  TokenType = "client"
)

type Role string

const (
	RoleSuperUser    Role = "superuser"
	RoleSupport      Role = "support"
	RoleStory        Role = "story"
	RoleFile         Role = "file"
	RoleNotification Role = "notification"
	RoleClient       Role = "client"
	RoleGuest        Role = "guest"
)

var AllUserRole = []Role{RoleSuperUser, RoleSupport, RoleStory, RoleFile, RoleNotification, RoleClient, RoleGuest}

var AdminRoles = []Role{RoleSuperUser, RoleSupport, RoleStory, RoleFile, RoleNotification}

func IsValidRole(role Role) bool {
	for _, r := range AllUserRole {
		if role == r {
			return true
		}
	}

	return false
}

func IsAdminRole(role Role) bool {
	for _, r := range AdminRoles {
		if role == r {
			return true
		}
	}

	return false
}
