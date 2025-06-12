package notificationapp

import (
	"context"
	"fmt"
	"log/slog"
	http2 "net/http"
	"os"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/app/notificationapp/delivery/http"
	postgres2 "github.com/syntaxfa/quick-connect/app/notificationapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/pkg/websocket"
)

type Application struct {
	cfg        Config
	trap       <-chan os.Signal
	logger     *slog.Logger
	httpServer http.Server
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal) Application {
	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		panic(tErr)
	}

	cfg.Hub.PingPeriod = (cfg.Hub.PongWait * 9) / 10

	redisAdapter := redis.New(cfg.Redis, logger)
	psAdapter := postgres.New(cfg.Postgres, logger)
	cache := cachemanager.New(redisAdapter)

	notificationVld := service.NewValidate(t)
	notificationRepo := postgres2.New(psAdapter)

	hub := service.NewHub(cfg.Hub, logger)
	upgrader := websocket.NewGorillaUpgrader(cfg.Websocket, checkOrigin(cfg.HTTPServer.Cors.AllowOrigins, logger))
	notificationSvc := service.New(cfg.Notification, notificationVld, cache, notificationRepo, logger, hub)

	handler := http.NewHandler(notificationSvc, t, upgrader)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), handler)

	return Application{
		cfg:        cfg,
		trap:       trap,
		logger:     logger,
		httpServer: httpServer,
	}
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

	a.logger.Info("notification app stopped")
}

func (a Application) Stop(ctx context.Context) bool {
	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup
		shutdownWg.Add(1)
		go a.StopHTTPServer(ctx, &shutdownWg)

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

func (a Application) StopHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	if sErr := a.httpServer.Stop(ctx); sErr != nil {
		a.logger.Error("http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}
}

func checkOrigin(allowOrigins []string, logger *slog.Logger) func(r *http2.Request) bool {
	return func(r *http2.Request) bool {
		_, _, _ = allowOrigins, logger, r

		return true
	}
}
