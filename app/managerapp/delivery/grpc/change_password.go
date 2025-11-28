package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (h Handler) UserChangePassword(ctx context.Context, req *userpb.UserChangePasswordRequest) (*empty.Empty, error) {
	userClaims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return &empty.Empty{}, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	if sErr := h.userSvc.ChangePassword(ctx, userClaims.UserID, userservice.ChangePasswordRequest{
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}); sErr != nil {
		return &empty.Empty{}, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return &empty.Empty{}, nil
}
