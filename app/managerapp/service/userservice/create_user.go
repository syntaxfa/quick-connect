package userservice

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) CreateUser(ctx context.Context, req UserCreateRequest) (UserCreateResponse, error) {
	op := "service.CreateUser"

	fmt.Printf("%+v\n\n\n", req)

	if vErr := s.vld.ValidateUserCreateRequest(req); vErr != nil {
		return UserCreateResponse{}, vErr
	}

	if uExist, ueErr := s.repo.IsExistUserByUsername(ctx, req.Username); ueErr != nil {
		return UserCreateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(ueErr).WithKind(richerror.KindUnexpected), s.logger)
	} else if uExist {
		return UserCreateResponse{}, richerror.New(op).WithKind(richerror.KindConflict).WithMessage(servermsg.MsgConflictUsername)
	}

	hashPass, hErr := HashPassword(req.Password)
	if hErr != nil {
		return UserCreateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(hErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	req.Password = hashPass
	req.ID = types.ID(ulid.Make().String())

	user, cErr := s.repo.CreateUser(ctx, req)
	if cErr != nil {
		return UserCreateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return UserCreateResponse{user}, nil
}
