package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Host string
	Port string
}

type Server struct {
	cfg        Config
	log        *slog.Logger
	GrpcServer *grpc.Server
}

func New(conf Config, log *slog.Logger) Server {
	grpc := grpc.NewServer()
	return Server{cfg: conf, log: log, GrpcServer: grpc}
}

func (s *Server) Start() error {
	s.log.Info("starting gRPC server at", slog.String("host", s.cfg.Host), slog.String("port", s.cfg.Port))

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port))
	if err != nil {
		s.log.Error("error at starting tcp listener", slog.String("error", err.Error()))
		return err
	}

	if err := s.GrpcServer.Serve(lis); err != nil {
		s.log.Error("error at serving grpc server", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("started gRPC successfully at", slog.String("host", s.cfg.Host), slog.String("port", s.cfg.Port))
	return nil
}

func (s *Server) Stop() {
	s.log.Info("gRPC server gracefully shutdown", slog.String("host", s.cfg.Host), slog.String("port", s.cfg.Port))
	s.GrpcServer.GracefulStop()
}
