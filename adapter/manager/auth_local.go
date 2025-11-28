package manager

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type AuthLocalAdapter struct {
	userSvc  *userservice.Service
	tokenSvc *tokenservice.Service
	t        *translation.Translate
	logger   *slog.Logger
}

func NewAuthLocalAdapter(userSvc *userservice.Service, tokenSvc *tokenservice.Service, t *translation.Translate,
	logger *slog.Logger) *AuthLocalAdapter {
	return &AuthLocalAdapter{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
		t:        t,
		logger:   logger,
	}
}

func (adl *AuthLocalAdapter) GetPublicKey(_ context.Context, _ *empty.Empty, _ ...grpc.CallOption) (*authpb.GetPublicKeyResponse, error) {
	return &authpb.GetPublicKeyResponse{PublicKey: adl.tokenSvc.GetPublicKey()}, nil
}

func (adl *AuthLocalAdapter) Login(ctx context.Context, req *authpb.LoginRequest, _ ...grpc.CallOption) (*authpb.LoginResponse, error) {
	resp, sErr := adl.userSvc.Login(ctx, convertLoginRequestToEntity(req))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, adl.t, adl.logger)
	}

	return convertLoginResponseToPB(resp), nil
}

func (adl *AuthLocalAdapter) TokenRefresh(ctx context.Context, req *authpb.TokenRefreshRequest,
	_ ...grpc.CallOption) (*authpb.TokenRefreshResponse, error) {
	resp, sErr := adl.userSvc.RefreshToken(ctx, req.GetRefreshToken())
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, adl.t, adl.logger)
	}

	return convertTokenGenerateResponseToPB(resp), nil
}
