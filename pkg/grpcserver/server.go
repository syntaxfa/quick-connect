package grpcserver

import (
	"fmt"
	"log/slog"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

type Server struct {
	cfg        Config
	log        *slog.Logger
	GrpcServer *grpc.Server
}

func New(conf Config, log *slog.Logger) Server {
	gr := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithMeterProvider(otel.GetMeterProvider()),
		)),
	)

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
