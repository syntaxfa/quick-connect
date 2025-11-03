package grpcserver

import (
	"context"
	"log/slog"
	"net"
	"strconv"

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

	//nolint:gocritic // Need separate slice for clarity
	allOpts := append(internalOpts, externalOpts...)

	gr := grpc.NewServer(allOpts...)

	// 4. (Recommendation) Register the reflection service on the server
	// This allows clients like grpcurl, Postman, etc., to discover
	// services and methods dynamically without needing the .proto files.
	// It's highly useful for debugging and development.
	reflection.Register(gr)

	return Server{cfg: conf, log: log, GrpcServer: gr}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.InfoContext(ctx, "starting gRPC server at", slog.String("host", s.cfg.Host), slog.Int("port", s.cfg.Port))

	lc := net.ListenConfig{}
	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))

	lis, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		s.log.ErrorContext(ctx, "error at starting tcp listener", slog.String("error", err.Error()))

		return err
	}

	if sErr := s.GrpcServer.Serve(lis); sErr != nil {
		s.log.ErrorContext(ctx, "error at serving grpcserver server", slog.String("error", sErr.Error()))

		return sErr
	}

	s.log.InfoContext(ctx, "started gRPC successfully at", slog.String("host", s.cfg.Host),
		slog.Int("port", s.cfg.Port))

	return nil
}

func (s *Server) Stop() {
	s.log.Info("gRPC server gracefully shutdown", slog.String("host", s.cfg.Host), slog.Int("port", s.cfg.Port))
	s.GrpcServer.GracefulStop()
}
