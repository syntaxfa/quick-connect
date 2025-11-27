package manager

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type AuthLocalAdapter struct {
	userSvc  *userservice.Service
	tokenSvc *tokenservice.Service
}

func NewAuthLocalAdapter(userSvc *userservice.Service, tokenSvc *tokenservice.Service) *AuthLocalAdapter {
	return &AuthLocalAdapter{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
	}
}

func (adl *AuthLocalAdapter) GetPublicKey(_ context.Context, _ *empty.Empty, _ ...grpc.CallOption) (*authpb.GetPublicKeyResponse, error) {
	return &authpb.GetPublicKeyResponse{PublicKey: adl.tokenSvc.GetPublicKey()}, nil
}
