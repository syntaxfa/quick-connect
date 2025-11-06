package userservice

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) RegisterGuestUser(ctx context.Context, req RegisterGuestUserRequest) (RegisterGuestUserResponse, error) {
	const op = "service.guest.RegisterGuestUser"

	if vErr := s.vld.ValidateRegisterGuestUserRequest(req); vErr != nil {
		return RegisterGuestUserResponse{}, vErr
	}

	userID := ulid.Make().String()
	userCreateReq := UserCreateRequest{
		ID:          types.ID(userID),
		Username:    userID,
		Password:    ulid.Make().String() + ulid.Make().String(),
		Fullname:    "",
		Email:       "",
		PhoneNumber: "",
		Roles:       []types.Role{types.RoleGuest},
	}

	if req.Fullname != "" {
		userCreateReq.Fullname = req.Fullname
	} else {
		userCreateReq.Fullname = userCreateReq.Username
	}

	if req.Email != "" {
		userCreateReq.Email = req.Email
	} else {
		userCreateReq.Email = fmt.Sprintf("%s@anomymous.none", userCreateReq.Username)
	}

	if req.PhoneNumber != "" {
		userCreateReq.PhoneNumber = req.PhoneNumber
	} else {
		userCreateReq.PhoneNumber = "None"
	}

	user, cErr := s.repo.CreateUser(ctx, userCreateReq)
	if cErr != nil {
		return RegisterGuestUserResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(cErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	qcToken, gErr := s.tokenSvc.GenerateGuestToken(ctx, user.ID)
	if gErr != nil {
		return RegisterGuestUserResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return RegisterGuestUserResponse{
		User:    user,
		QCToken: qcToken,
	}, nil
}
