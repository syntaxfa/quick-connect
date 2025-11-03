package jwtvalidator

import (
	"crypto"
	"crypto/ed25519"
	"encoding/hex"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

type Validator struct {
	publicKeyString string
	publicKey       crypto.PublicKey
	logger          *slog.Logger
}

func New(publicKeyString string, logger *slog.Logger) *Validator {
	return &Validator{
		publicKeyString: publicKeyString,
		logger:          logger,
	}
}

func (v *Validator) ValidateToken(tokenString string) (*types.UserClaims, error) {
	op := "token.service.ValidateToken"

	publicKeyBytes, dErr := hex.DecodeString(v.publicKeyString)
	if dErr != nil {
		return nil, errlog.ErrLog(richerror.New(op).WithWrapError(dErr).WithMessage("invalid public key").
			WithKind(richerror.KindUnAuthorized), v.logger)
	}

	v.publicKey = ed25519.PublicKey(publicKeyBytes)

	token, err := jwt.ParseWithClaims(tokenString, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate algorithm signature.
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidTokenAlgorithm).
				WithKind(richerror.KindUnAuthorized), v.logger)
		}

		return v.publicKey, nil
	})
	if err != nil {
		return nil, errlog.ErrLog(richerror.New(op).WithMessage(servermsg.MsgInvalidToken).WithWrapError(err).
			WithKind(richerror.KindUnAuthorized), v.logger)
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
