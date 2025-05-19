package main

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/example/grpc/proto/pub"
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type greeterServer struct {
	pub.UnimplementedGreeterServer
}

func SayHello(ctx context.Context, in *pub.HelloRequest) (*pub.HelloReply, error) {
	return &pub.HelloReply{Message: "Hello, World! "}, nil
}

func main() {
	log := logger.New(logger.Config{
		FilePath:         "logs.json",
		UseLocalTime:     false,
		FileMaxSizeInMB:  1,
		FileMaxAgeInDays: 10,
		MaxBackup:        0,
		Compress:         false,
	}, nil, true, "example")

	server := grpcserver.New(grpcserver.Config{
		Host: "localhost",
		Port: 50051,
	}, log)

	go server.Start()
	pub.RegisterGreeterServer(server.GrpcServer, greeterServer{})

	time.Sleep(time.Second * 5)

	server.Stop()
}
