package jwtvalidator

import (
	"crypto"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

type Validator struct {
	publicKey crypto.PublicKey
	logger    *slog.Logger
}

func New(publicKey crypto.PublicKey, logger *slog.Logger) *Validator {
	return &Validator{
		publicKey: publicKey,
		logger:    logger,
	}
}

func (v *Validator) ValidateToken(tokenString string) (*types.UserClaims, error) {
	op := "token.service.ValidateToken"

	token, err := jwt.ParseWithClaims(tokenString, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate algorithm signature.
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidTokenAlgorithm).WithKind(richerror.KindUnAuthorized), v.logger)
		}

		return v.publicKey, nil
	})
	if err != nil {
		return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidToken).WithWrapError(err).WithKind(richerror.KindUnAuthorized), v.logger)
	}

	if !token.Valid {
		return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidToken).WithKind(richerror.KindUnAuthorized), v.logger)
	}

	claims, ok := token.Claims.(*types.UserClaims)
	if !ok {
		return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidToken).WithKind(richerror.KindUnAuthorized), v.logger)
	}

	return claims, nil
}
