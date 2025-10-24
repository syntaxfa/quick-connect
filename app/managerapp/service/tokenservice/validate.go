package tokenservice

import (
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) ValidateToken(tokenString string) (*types.UserClaims, error) {
	return s.validator.ValidateToken(tokenString)
}
