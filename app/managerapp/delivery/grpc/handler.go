package grpc

import (
	"context"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"log/slog"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/tokenpb"
)

type Handler struct {
	tokenpb.UnimplementedTokenServiceServer
	logger   *slog.Logger
	tokenSvc tokenservice.Service
}

func NewHandler(logger *slog.Logger, tokenSvc tokenservice.Service) Handler {
	return Handler{
		logger:   logger,
		tokenSvc: tokenSvc,
	}
}

func (h Handler) GetPublicKey(_ context.Context, _ *empty.Empty) (*tokenpb.GetPublicKeyResponse, error) {
	return &tokenpb.GetPublicKeyResponse{PublicKey: h.tokenSvc.GetPublicKey()}, nil
}
