package tokenservice

import (
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) RefreshTokens(refreshToken string) (*TokenGenerateResponse, error) {
	op := "auth.service.RefreshTokens"

	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, richerror.New(op).WithWrapError(err).WithMessage(servermsg.MsgInvalidToken).
			WithKind(richerror.KindUnAuthorized)
	}

	if claims.TokenType != types.TokenTypeRefresh {
		return nil, richerror.New(op).WithMessage(servermsg.MsgInvalidToken).WithKind(richerror.KindUnAuthorized)
	}

	return s.GenerateTokenPair(claims.UserID, claims.Role)
}
