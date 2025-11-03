package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/protobuf/example/golang/examplepb"
)

type Server struct {
	grpcServer grpcserver.Server
	handler    Handler
}

func New(server grpcserver.Server, handler Handler) Server {
	return Server{
		grpcServer: server,
		handler:    handler,
	}
}

func (s Server) Start() error {
	examplepb.RegisterCommentServiceServer(s.grpcServer.GrpcServer, s.handler)

	return s.grpcServer.Start(context.Background())
}

func (s Server) Stop() {
	s.grpcServer.Stop()
}
