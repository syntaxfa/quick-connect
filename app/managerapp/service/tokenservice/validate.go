package tokenservice

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) ValidateToken(tokenString string) (*types.UserClaims, error) {
	op := "token.service.ValidateToken"

	token, err := jwt.ParseWithClaims(tokenString, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate algorithm signature.
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidTokenAlgorithm), s.logger)
		}

		return s.cfg.publicKey, nil
	})
	if err != nil {
		return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidToken).WithWrapError(err), s.logger)
	}

	if !token.Valid {
		return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidToken), s.logger)
	}

	claims, ok := token.Claims.(*types.UserClaims)
	if !ok {
		return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidToken), s.logger)
	}

	return claims, nil
}
