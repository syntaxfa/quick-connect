package tokenservice

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) GenerateTokenPair(userID types.ID, roles []types.Role) (*TokenGenerateResponse, error) {
	op := "auth.service.GenerateTokenPair"

	accessToken, gaErr := s.generateToken(userID, roles, types.TokenTypeAccess, s.cfg.AccessExpiry, s.cfg.AccessAudience)
	if gaErr != nil {
		return nil, errlog.ErrLog(richerror.New(op).WithWrapError(gaErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	refreshToken, grErr := s.generateToken(userID, roles, types.TokenTypeRefresh, s.cfg.RefreshExpiry, s.cfg.RefreshAudience)
	if grErr != nil {
		return nil, errlog.ErrLog(richerror.New(op).WithWrapError(grErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return &TokenGenerateResponse{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresIn: int32(s.cfg.AccessExpiry.Seconds()),
		RefreshExpireIn: int32(s.cfg.RefreshExpiry.Seconds()),
	}, nil
}

func (s Service) generateToken(userID types.ID, roles []types.Role, tokenType types.TokenType, expiry time.Duration,
	audience string) (string, error) {
	op := "auth.service.generateToken"

	now := time.Now().UTC()

	// unique token id
	tokenID, uErr := uuid.NewRandom()
	if uErr != nil {
		return "", errlog.ErrLog(richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	claims := types.UserClaims{
		UserID:    userID,
		Roles:     roles,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.cfg.Issuer,
			Subject:   fmt.Sprintf("%s token for user %s", tokenType, userID),
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signedToken, err := token.SignedString(s.cfg.privateKey)
	if err != nil {
		return "", errlog.ErrLog(richerror.New(op).WithMessage(err.Error()).WithWrapError(err).
			WithKind(richerror.KindUnAuthorized), s.logger)
	}

	return signedToken, nil
}

func (s Service) GenerateGuestToken(ctx context.Context, userID types.ID) (string, error) {
	const op = "auth.service.generate.GenerateGuestToken"

	guestToken, gErr := s.generateToken(userID, []types.Role{types.RoleGuest}, types.TokenTypeGuest, s.cfg.GuestExpiry,
		s.cfg.GuestAudience)
	if gErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return guestToken, nil
}

func (s Service) GenerateClientToken(ctx context.Context, userID types.ID, roles []types.Role) (string, error) {
	const op = "auth.service.generate.GenerateClientToken"

	guestToken, gErr := s.generateToken(userID, roles, types.TokenTypeClient, s.cfg.ClientExpiry,
		s.cfg.GuestAudience)
	if gErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return guestToken, nil
}
