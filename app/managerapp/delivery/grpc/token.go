package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
)

func (h Handler) GetPublicKey(_ context.Context, _ *empty.Empty) (*authpb.GetPublicKeyResponse, error) {
	return &authpb.GetPublicKeyResponse{PublicKey: h.tokenSvc.GetPublicKey()}, nil
}
