package manager

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"google.golang.org/grpc"
)

type AuthAdapter struct {
	client authpb.AuthServiceClient
}

func NewAuthAdapter(conn *grpc.ClientConn) *AuthAdapter {
	return &AuthAdapter{
		client: authpb.NewAuthServiceClient(conn),
	}
}

func (tc AuthAdapter) GetPublicKey(ctx context.Context) (*authpb.GetPublicKeyResponse, error) {
	return tc.client.GetPublicKey(ctx, &empty.Empty{})
}

func (tc AuthAdapter) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return tc.client.Login(ctx, req)
}
