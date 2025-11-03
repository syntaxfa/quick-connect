package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (h Handler) GetPublicKey(_ context.Context, _ *empty.Empty) (*authpb.GetPublicKeyResponse, error) {
	return &authpb.GetPublicKeyResponse{PublicKey: h.tokenSvc.GetPublicKey()}, nil
}
