package grpc

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
)

type Handler struct {
	authpb.UnimplementedAuthServiceServer
	userpb.UnimplementedUserServiceServer
	logger   *slog.Logger
	tokenSvc tokenservice.Service
	userSvc  userservice.Service
	t        *translation.Translate
}

func NewHandler(logger *slog.Logger, tokenSvc tokenservice.Service, userSvc userservice.Service, t *translation.Translate) Handler {
	return Handler{
		logger:   logger,
		tokenSvc: tokenSvc,
		userSvc:  userSvc,
		t:        t,
	}
}
