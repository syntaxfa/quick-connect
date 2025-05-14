package userservice

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func (s Service) Login(ctx context.Context, req UserLoginRequest) (UserLoginResponse, error) {
	const op = "service.Login"

	if vErr := s.vld.ValidateLoginRequest(req); vErr != nil {
		return UserLoginResponse{}, vErr
	}

	if exists, iErr := s.repo.IsExistUserByUsername(ctx, req.Username); iErr != nil {
		return UserLoginResponse{}, richerror.New(op).WithWrapError(iErr).WithKind(richerror.KindUnexpected)
	} else if !exists {
		return UserLoginResponse{}, richerror.New(op).WithKind(richerror.KindNotFound)
	}

	user, gErr := s.repo.GetUserByUsername(ctx, req.Username)
	if gErr != nil {
		return UserLoginResponse{}, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected)
	}

	if !VerifyPassword(user.HashedPassword, req.Password) {
		return UserLoginResponse{}, richerror.New(op).WithKind(richerror.KindNotFound)
	}

	token, gtErr := s.tokenSvc.GenerateTokenPair(user.ID, user.Role)
	if gtErr != nil {
		return UserLoginResponse{}, richerror.New(op).WithWrapError(gtErr).WithKind(richerror.KindUnexpected)
	}

	return UserLoginResponse{
		User:  user,
		Token: *token,
	}, nil
}
