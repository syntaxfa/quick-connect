package types

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserID    ID        `json:"user_id"`
	Role      Role      `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Role uint8

const (
	RoleSuperUser = iota + 1
	RoleAdmin
)
