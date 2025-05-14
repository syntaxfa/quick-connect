package tokenservice

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/types"
)

type CustomClaims struct {
	UserID    types.ID         `json:"user_id"`
	Role      userservice.Role `json:"role"`
	TokenType TokenType        `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)
