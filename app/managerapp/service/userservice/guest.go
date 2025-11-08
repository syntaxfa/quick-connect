package userservice

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
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
	}

	if req.Email != "" {
		userCreateReq.Email = req.Email
	}

	if req.PhoneNumber != "" {
		userCreateReq.PhoneNumber = req.PhoneNumber
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

func (s Service) UpdateGuestUser(ctx context.Context, userID types.ID, req UpdateGuestUserRequest) (UpdateGuestUserResponse, error) {
	const op = "service.guest.UpdateGuestUser"

	if vErr := s.vld.ValidateUpdateGuestUserRequest(req); vErr != nil {
		return UpdateGuestUserResponse{}, vErr
	}

	exists, existErr := s.repo.IsExistUserByID(ctx, userID)
	if existErr != nil {
		return UpdateGuestUserResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(existErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		return UpdateGuestUserResponse{}, richerror.New(op).WithMessage(servermsg.MsgRecordNotFound).WithKind(richerror.KindNotFound)
	}

	user, gErr := s.repo.GetUserByID(ctx, userID)
	if gErr != nil {
		return UpdateGuestUserResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	updateUser := UserUpdateFromSuperuserRequest{Username: user.Username, Roles: user.Roles}

	if req.Fullname != "" {
		updateUser.Fullname = req.Fullname
	} else {
		updateUser.Fullname = user.Fullname
	}

	if req.Email != "" {
		updateUser.Email = req.Email
	} else {
		updateUser.Email = user.Email
	}

	if req.PhoneNumber != "" {
		updateUser.PhoneNumber = req.PhoneNumber
	} else {
		updateUser.PhoneNumber = user.PhoneNumber
	}

	if uErr := s.repo.UpdateUser(ctx, userID, updateUser); uErr != nil {
		return UpdateGuestUserResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(uErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return UpdateGuestUserResponse{
		User{
			ID:             user.ID,
			Username:       updateUser.Username,
			HashedPassword: "",
			Fullname:       updateUser.Fullname,
			Email:          updateUser.Email,
			PhoneNumber:    updateUser.PhoneNumber,
			Avatar:         user.Avatar,
			Roles:          user.Roles,
			LastOnlineAt:   user.LastOnlineAt,
		},
	}, nil
}
