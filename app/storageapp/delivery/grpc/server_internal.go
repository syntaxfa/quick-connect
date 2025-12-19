package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
)

type InternalServer struct {
	server  grpcserver.Server
	handler InternalHandler
}

func New(server grpcserver.Server, handler InternalHandler) InternalServer {
	return InternalServer{
		server:  server,
		handler: handler,
	}
}

func (s InternalServer) Start(ctx context.Context) error {
	storagepb.RegisterStorageInternalServiceServer(s.server.GrpcServer, s.handler)

	return s.server.Start(ctx)
}

func (s InternalServer) Stop() {
	s.server.Stop()
}
