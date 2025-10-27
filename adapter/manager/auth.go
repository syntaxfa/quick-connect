package manager

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/tokenpb"
	"google.golang.org/grpc"
)

type TokenAdapter struct {
	client tokenpb.TokenServiceClient
}

func NewTokenAdapter(conn *grpc.ClientConn) *TokenAdapter {
	return &TokenAdapter{
		client: tokenpb.NewTokenServiceClient(conn),
	}
}

func (tc TokenAdapter) GetPublicKey(ctx context.Context) (*tokenpb.GetPublicKeyResponse, error) {
	return tc.client.GetPublicKey(ctx, &empty.Empty{})
}
