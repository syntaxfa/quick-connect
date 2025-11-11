package userservice

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) RefreshToken(ctx context.Context, refreshToken string) (*tokenservice.TokenGenerateResponse, error) {
	const op = "userservice.refresh.RefreshToken"

	claims, tvErr := s.tokenSvc.ValidateToken(refreshToken)
	if tvErr != nil {
		return nil, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(tvErr).WithMessage(servermsg.MsgInvalidToken).
			WithKind(richerror.KindUnAuthorized), s.logger)
	}

	if claims.TokenType != types.TokenTypeRefresh {
		return nil, errlog.ErrContext(ctx, richerror.New(op).WithMessage(servermsg.MsgInvalidToken).
			WithKind(richerror.KindUnAuthorized), s.logger)
	}

	exists, existErr := s.repo.IsExistUserByID(ctx, claims.UserID)
	if existErr != nil {
		return nil, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(existErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		return nil, errlog.ErrContext(ctx, richerror.New(op).WithMessage(servermsg.MsgInvalidToken).
			WithKind(richerror.KindUnAuthorized), s.logger)
	}

	user, guErr := s.repo.GetUserByID(ctx, claims.UserID)
	if guErr != nil {
		return nil, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(guErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return s.tokenSvc.GenerateTokenPair(user.ID, user.Roles)
}
