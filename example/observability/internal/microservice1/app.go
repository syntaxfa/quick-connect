package microservice1

import (
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice1/delivery/http"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice1/repository"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice1/service"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"log/slog"
	"os"
)

type Application struct {
	cfg        Config
	handler    http.Handler
	trap       <-chan os.Signal
	logger     *slog.Logger
	httpServer http.Server
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal) Application {
	ps := postgres.New(cfg.Postgres)
	repo := repository.New(ps)
	svc := service.New(repo)

	handler := http.NewHandler(svc)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), handler, cfg.Observability.Core.ServiceName)

	return Application{
		cfg:        cfg,
		handler:    handler,
		trap:       trap,
		logger:     logger,
		httpServer: httpServer,
	}
}

func (a Application) Start() {
	go func() {
		a.logger.Info(fmt.Sprintf("http server started on %d", a.cfg.HTTPServer.Port))

		if sErr := a.httpServer.Start(); sErr != nil {
			a.logger.Error(fmt.Sprintf("error in http server on %d", a.cfg.HTTPServer.Port), slog.String("error", sErr.Error()))
		}
		a.logger.Info(fmt.Sprintf("http server stopped %d", a.cfg.HTTPServer.Port))
	}()

	<-a.trap

	a.logger.Info("shutdown signal received")

	ctx, cancelFunc := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancelFunc()

	a.Stop(ctx)
}

func (a Application) Stop(ctx context.Context) {
	if sErr := a.httpServer.Stop(ctx); sErr != nil {
		a.logger.Error("http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}

	a.logger.Info("http server gracefully shutdown")
}
