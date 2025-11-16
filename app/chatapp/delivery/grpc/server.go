package grpc

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
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
	conversationpb.RegisterConversationServiceServer(s.server.GrpcServer, s.handler)

	return s.server.Start(context.Background())
}

func (s Server) Stop() {
	s.server.Stop()
}
