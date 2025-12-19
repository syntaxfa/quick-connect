package storyapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/syntaxfa/quick-connect/app/storyapp/delivery/http"
	"github.com/syntaxfa/quick-connect/app/storyapp/service"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Application struct {
	cfg        Config
	httpServer http.Server
	logger     *slog.Logger
	trap       <-chan os.Signal
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal, t *translation.Translate) (Application, service.Service) {
	svc := service.New()

	handler := http.NewHandler(svc, t, logger)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), handler)

	return Application{
		cfg:        cfg,
		httpServer: httpServer,
		logger:     logger,
		trap:       trap,
	}, svc
}

func (a Application) Start() {
	httpServerChan := make(chan error, 1)

	go func() {
		a.logger.Info(fmt.Sprintf("http server started on %d", a.cfg.HTTPServer.Port))

		if sErr := a.httpServer.Start(); sErr != nil {
			httpServerChan <- sErr
		}
	}()

	select {
	case err := <-httpServerChan:
		a.logger.Error(fmt.Sprintf("error in http server on %d", a.cfg.HTTPServer.Port), slog.String("error", err.Error()))
	case <-a.trap:
		a.logger.Info("received http server shutdown signal!!!")
	}

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	if a.Stop(shutdownTimeoutCtx) {
		a.logger.Info("servers shutdown gracefully")
	} else {
		a.logger.Warn("shutdown timed out, existing application")
	}

	a.logger.Info("story app stopped")
}

func (a Application) Stop(ctx context.Context) bool {
	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup

		shutdownWg.Add(1)
		go a.stopHTTPServer(ctx, &shutdownWg)

		shutdownWg.Wait()
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
		return true
	case <-ctx.Done():
		return false
	}
}

func (a Application) stopHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if sErr := a.httpServer.Stop(ctx); sErr != nil {
		a.logger.ErrorContext(ctx, "http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}
}
