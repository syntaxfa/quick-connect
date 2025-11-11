package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
)

func (h Handler) TokenRefresh(ctx context.Context, req *authpb.TokenRefreshRequest) (*authpb.TokenRefreshResponse, error) {
	resp, sErr := h.userSvc.RefreshToken(ctx, req.GetRefreshToken())
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return &authpb.TokenRefreshResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresIn:  resp.AccessExpiresIn,
		RefreshExpiresIn: resp.RefreshExpireIn,
	}, nil
}
