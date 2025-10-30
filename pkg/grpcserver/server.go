package grpcserver

import (
	"fmt"
	"log/slog"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg        Config
	log        *slog.Logger
	GrpcServer *grpc.Server
}

func New(conf Config, log *slog.Logger, externalOpts ...grpc.ServerOption) Server {
	internalOpts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithMeterProvider(otel.GetMeterProvider()),
		)),
	}

	allOpts := append(internalOpts, externalOpts...)

	gr := grpc.NewServer(allOpts...)

	// 4. (Recommendation) Register the reflection service on the server
	// This allows clients like grpcurl, Postman, etc., to discover
	// services and methods dynamically without needing the .proto files.
	// It's highly useful for debugging and development.
	reflection.Register(gr)

	return Server{cfg: conf, log: log, GrpcServer: gr}
}

func (s *Server) Start() error {
	s.log.Info("starting gRPC server at", slog.String("host", s.cfg.Host), slog.Int("port", s.cfg.Port))

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
	if err != nil {
		s.log.Error("error at starting tcp listener", slog.String("error", err.Error()))

		return err
	}

	if err := s.GrpcServer.Serve(lis); err != nil {
		s.log.Error("error at serving grpcserver server", slog.String("error", err.Error()))

		return err
	}

	s.log.Info("started gRPC successfully at", slog.String("host", s.cfg.Host), slog.Int("port", s.cfg.Port))

	return nil
}

func (s *Server) Stop() {
	s.log.Info("gRPC server gracefully shutdown", slog.String("host", s.cfg.Host), slog.Int("port", s.cfg.Port))
	s.GrpcServer.GracefulStop()
}
