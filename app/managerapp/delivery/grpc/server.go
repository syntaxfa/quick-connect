package grpc

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
)

type Server struct {
	server  grpcserver.Server
	handler Handler
	logger  *slog.Logger
}

func New(server grpcserver.Server, handler Handler, logger *slog.Logger) Server {
	return Server{
		server:  server,
		handler: handler,
		logger:  logger,
	}
}

func (s Server) Start() error {
	authpb.RegisterAuthServiceServer(s.server.GrpcServer, s.handler)
	userpb.RegisterUserServiceServer(s.server.GrpcServer, s.handler)

	return s.server.Start()
}

func (s Server) Stop() {
	s.server.Stop()
}
