package manager

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type UserLocalAdapter struct {
	userSvc *userservice.Service
	t       *translation.Translate
	logger  *slog.Logger
}

func NewUserLocalAdapter(userSvc *userservice.Service, t *translation.Translate, logger *slog.Logger) *UserLocalAdapter {
	return &UserLocalAdapter{
		userSvc: userSvc,
		t:       t,
		logger:  logger,
	}
}

func (udl *UserLocalAdapter) UserProfile(ctx context.Context, _ *empty.Empty, _ ...grpc.CallOption) (*userpb.User, error) {
	claims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	resp, sErr := udl.userSvc.UserProfile(ctx, claims.UserID)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, udl.t, udl.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (udl *UserLocalAdapter) UserUpdateFromOwn(ctx context.Context, req *userpb.UserUpdateFromOwnRequest,
	_ ...grpc.CallOption) (*userpb.User, error) {
	claims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	resp, sErr := udl.userSvc.UserUpdateFromOwn(ctx, claims.UserID, convertUserUpdateFromOwnToEntity(req))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, udl.t, udl.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (udl *UserLocalAdapter) UserList(ctx context.Context, req *userpb.UserListRequest,
	_ ...grpc.CallOption) (*userpb.UserListResponse, error) {
	resp, sErr := udl.userSvc.UserList(ctx, convertUserListRequestToEntity(req))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, udl.t, udl.logger)
	}

	return convertUserListResponseToPB(resp), nil
}

func (udl *UserLocalAdapter) UserDelete(ctx context.Context, req *userpb.UserDeleteRequest, _ ...grpc.CallOption) (*empty.Empty, error) {
	if sErr := udl.userSvc.UserDelete(ctx, types.ID(req.GetUserId())); sErr != nil {
		return &empty.Empty{}, servermsg.GRPCMsg(sErr, udl.t, udl.logger)
	}

	return &empty.Empty{}, nil
}

func (udl *UserLocalAdapter) UserDetail(ctx context.Context, req *userpb.UserDetailRequest, _ ...grpc.CallOption) (*userpb.User, error) {
	resp, sErr := udl.userSvc.UserProfile(ctx, types.ID(req.GetUserId()))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, udl.t, udl.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (udl *UserLocalAdapter) UserUpdateFromSuperuser(ctx context.Context, req *userpb.UserUpdateFromSuperUserRequest,
	_ ...grpc.CallOption) (*userpb.User, error) {
	resp, sErr := udl.userSvc.UserUpdateFromSuperuser(ctx, types.ID(req.GetUserId()), convertUserUpdateFromSuperuserToEntity(req))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, udl.t, udl.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (udl *UserLocalAdapter) CreateUser(ctx context.Context, req *userpb.CreateUserRequest, _ ...grpc.CallOption) (*userpb.User, error) {
	resp, sErr := udl.userSvc.CreateUser(ctx, convertCreateUserRequestToEntity(req))

	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, udl.t, udl.logger)
	}

	return convertUserToPB(resp.User), nil
}

func (udl *UserLocalAdapter) UserChangePassword(ctx context.Context, req *userpb.UserChangePasswordRequest,
	_ ...grpc.CallOption) (*empty.Empty, error) {
	userClaims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return &empty.Empty{}, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	if sErr := udl.userSvc.ChangePassword(ctx, userClaims.UserID, userservice.ChangePasswordRequest{
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}); sErr != nil {
		return &empty.Empty{}, servermsg.GRPCMsg(sErr, udl.t, udl.logger)
	}

	return &empty.Empty{}, nil
}
