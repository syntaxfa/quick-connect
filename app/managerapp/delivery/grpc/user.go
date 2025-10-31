package grpc

import (
	"context"
	"fmt"

	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"github.com/syntaxfa/quick-connect/types"
)

func (h Handler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	resp, sErr := h.userSvc.CreateUser(ctx, convertCreateUserRequestToEntity(req))

	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (h Handler) UserDetail(ctx context.Context, req *userpb.UserDetailRequest) (*userpb.User, error) {
	fmt.Println("UserDetail server")
	resp, sErr := h.userSvc.UserProfile(ctx, types.ID(req.UserId))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserToPB(resp.User), nil
}
