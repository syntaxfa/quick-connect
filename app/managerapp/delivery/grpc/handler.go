package grpc

import (
	"context"
	"log/slog"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/tokenpb"
)

type Handler struct {
	tokenpb.UnimplementedTokenServiceServer
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) Handler {
	return Handler{
		logger: logger,
	}
}

func (h Handler) GetPublicKey(ctx context.Context, empty *empty.Empty) (*tokenpb.GetPublicKeyResponse, error) {
	return nil, nil
}
