package manager

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"google.golang.org/grpc"
)

type TokenAdapter struct {
	tokenClient authpb.TokenServiceClient
}

func NewTokenAdapter(conn *grpc.ClientConn) *TokenAdapter {
	return &TokenAdapter{
		tokenClient: authpb.NewTokenServiceClient(conn),
	}
}

func (tc TokenAdapter) GetPublicKey(ctx context.Context) (*authpb.GetPublicKeyResponse, error) {
	return tc.tokenClient.GetPublicKey(ctx, &empty.Empty{})
}
