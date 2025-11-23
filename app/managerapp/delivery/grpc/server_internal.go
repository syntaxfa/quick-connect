package grpc

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userinternalpb"
)

type ServerInternal struct {
	server  grpcserver.Server
	handler HandlerInternal
	logger  *slog.Logger
}

func NewServerInternal(server grpcserver.Server, handler HandlerInternal, logger *slog.Logger) ServerInternal {
	return ServerInternal{
		server:  server,
		handler: handler,
		logger:  logger,
	}
}

func (s ServerInternal) Start() error {
	userinternalpb.RegisterUserInternalServiceServer(s.server.GrpcServer, s.handler)

	return s.server.Start(context.Background())
}

func (s ServerInternal) Stop() {
	s.server.Stop()
}
