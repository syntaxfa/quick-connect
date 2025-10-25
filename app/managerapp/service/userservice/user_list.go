package userservice

import (
	"context"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func (s Service) UserList(ctx context.Context, req ListUserRequest) (ListUserResponse, error) {
	const op = "service.user_list.UserList"

	if vErr := s.vld.ValidateListUserRequest(req); vErr != nil {
		return ListUserResponse{}, vErr
	}

	if bErr := req.Paginated.BasicValidation(); bErr != nil {
		return ListUserResponse{}, richerror.New(op).WithKind(richerror.KindBadRequest)
	}

	users, paginateRes, uErr := s.repo.GetUserList(ctx, req.Paginated, req.Username)
	if uErr != nil {
		return ListUserResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return ListUserResponse{
		Results:  users,
		Paginate: paginateRes,
	}, nil
}
