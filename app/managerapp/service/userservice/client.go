package userservice

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/randomly"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) IdentifyClient(ctx context.Context, req IdentifyClientRequest) (IdentifyClientResponse, error) {
	const op = "service.client.IdentifyClient"

	if vErr := s.vld.ValidateIdentifyClientRequest(req); vErr != nil {
		return IdentifyClientResponse{}, vErr
	}

	userID, exErr := s.GetUserIDFromExternalUserID(ctx, req.ExternalUserID)
	if exErr != nil {
		return IdentifyClientResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(exErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	exists, existErr := s.repo.IsExistUserByID(ctx, userID)
	if existErr != nil {
		return IdentifyClientResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(existErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		password, passErr := randomly.GeneratePassword(s.cfg.PasswordDefaultLength)
		if passErr != nil {
			return IdentifyClientResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(passErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}

		createReq := UserCreateRequest{
			ID:          userID,
			Username:    string(userID),
			Password:    password,
			Fullname:    req.Fullname,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Roles:       []types.Role{types.RoleClient},
		}

		user, createErr := s.repo.CreateUser(ctx, createReq)
		if createErr != nil {
			return IdentifyClientResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(createErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}

		qcToken, tokenErr := s.tokenSvc.GenerateClientToken(ctx, user.ID, user.Roles)
		if tokenErr != nil {
			return IdentifyClientResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(tokenErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}

		return IdentifyClientResponse{
			User:    user,
			QCToken: qcToken,
		}, nil
	}

	user, getErr := s.repo.GetUserByID(ctx, userID)
	if getErr != nil {
		return IdentifyClientResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(getErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	updateReq := UserUpdateFromSuperuserRequest{
		Username:    user.Username,
		Fullname:    user.Fullname,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Roles:       user.Roles,
	}

	if req.Fullname != "" {
		updateReq.Fullname = req.Fullname
		user.Fullname = req.Fullname
	}

	if req.Email != "" {
		updateReq.Email = req.Email
		user.Email = req.Email
	}

	if req.PhoneNumber != "" {
		updateReq.PhoneNumber = req.PhoneNumber
		user.PhoneNumber = req.PhoneNumber
	}

	if upErr := s.repo.UpdateUser(ctx, userID, updateReq); upErr != nil {
		return IdentifyClientResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(upErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	qcToken, tokenErr := s.tokenSvc.GenerateClientToken(ctx, userID, user.Roles)
	if tokenErr != nil {
		return IdentifyClientResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(tokenErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return IdentifyClientResponse{
		User:    user,
		QCToken: qcToken,
	}, nil
}
