package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
)

func (h Handler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	resp, sErr := h.userSvc.CreateUser(ctx, convertCreateUserRequestToEntity(req))

	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserToPB(resp.User), nil
}
