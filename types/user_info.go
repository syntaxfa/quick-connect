package types

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserID    ID        `json:"user_id"`
	Roles     []Role    `json:"roles"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Role string

const (
	RoleSuperUser    Role = "superuser"
	RoleSupport      Role = "support"
	RoleStory        Role = "story"
	RoleFile         Role = "file"
	RoleNotification Role = "notification"
)

func IsValidRole(role Role) bool {
	if role == RoleSuperUser || role == RoleSupport || role == RoleStory || role == RoleFile || role == RoleNotification {
		return true
	}

	return false
}
