package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	resp, sErr := h.userSvc.CreateUser(ctx, convertCreateUserRequestToEntity(req))

	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (h Handler) UserDetail(ctx context.Context, req *userpb.UserDetailRequest) (*userpb.User, error) {
	resp, sErr := h.userSvc.UserProfile(ctx, types.ID(req.UserId))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (h Handler) UserDelete(ctx context.Context, req *userpb.UserDeleteRequest) (*empty.Empty, error) {
	if sErr := h.userSvc.UserDelete(ctx, types.ID(req.UserId)); sErr != nil {
		return &empty.Empty{}, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return &empty.Empty{}, nil
}

func (h Handler) UserUpdateFromSuperuser(ctx context.Context, req *userpb.UserUpdateFromSuperUserRequest) (*userpb.User, error) {
	resp, sErr := h.userSvc.UserUpdateFromSuperuser(ctx, types.ID(req.UserId), convertUserUpdateFromSuperuserToEntity(req))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (h Handler) UserUpdateFromOwn(ctx context.Context, req *userpb.UserUpdateFromOwnRequest) (*userpb.User, error) {
	userClaims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	resp, sErr := h.userSvc.UserUpdateFromOwn(ctx, userClaims.UserID, convertUserUpdateFromOwnToEntity(req))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (h Handler) UserList(ctx context.Context, req *userpb.UserListRequest) (*userpb.UserListResponse, error) {
	resp, sErr := h.userSvc.UserList(ctx, convertUserListRequestToEntity(req))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserListResponseToPB(resp), nil
}

func (h Handler) UserProfile(ctx context.Context, _ *empty.Empty) (*userpb.User, error) {
	userClaims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	resp, sErr := h.userSvc.UserProfile(ctx, userClaims.UserID)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserToPB(resp.User), nil
}
