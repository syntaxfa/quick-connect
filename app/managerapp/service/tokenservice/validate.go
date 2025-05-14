package tokenservice

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

func (s Service) ValidateToken(tokenString string) (*CustomClaims, error) {
	op := "token.service.ValidateToken"

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate algorithm signature.
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, richerror.New(op).WithMessage(servermsg.MsgInvalidTokenAlgorithm)
		}

		return s.cfg.publicKey, nil
	})
	if err != nil {
		return nil, richerror.New(op).WithMessage(servermsg.MsgInvalidToken).WithWrapError(err)
	}

	if !token.Valid {
		return nil, richerror.New(op).WithMessage(servermsg.MsgInvalidToken)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, richerror.New(op).WithMessage(servermsg.MsgInvalidToken)
	}

	return claims, nil
}
