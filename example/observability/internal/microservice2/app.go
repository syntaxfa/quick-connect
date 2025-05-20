package microservice2

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice2/delivery/grpc"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice2/repository"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice2/service"
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
)

type Application struct {
	cfg        Config
	grpcServer grpc.Server
	logger     *slog.Logger
	trap       <-chan os.Signal
}

func New(cfg Config) Application {
	return Application{
		cfg: cfg,
	}
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal) Application {
	ps := postgres.New(cfg.Postgres, logger)
	repo := repository.New(ps)
	svc := service.New(repo)
	handler := grpc.NewHandler(svc)
	grpcServer := grpc.New(grpcserver.New(cfg.GRPCServer, logger), handler)

	return Application{
		cfg:        cfg,
		grpcServer: grpcServer,
		logger:     logger,
		trap:       trap,
	}
}

func (a Application) Start() {
	go func() {
		a.logger.Info(fmt.Sprintf("groc server started on %d", a.cfg.GRPCServer.Port))

		if sErr := a.grpcServer.Start(); sErr != nil {
			a.logger.Error(fmt.Sprintf("error in grpc server on %d", a.cfg.GRPCServer.Port), slog.String("error", sErr.Error()))
		}
		a.logger.Info(fmt.Sprintf("groc server stopped %d", a.cfg.GRPCServer.Port))
	}()

	<-a.trap

	a.logger.Info("shutdown signal received")

	a.grpcServer.Stop()
}
