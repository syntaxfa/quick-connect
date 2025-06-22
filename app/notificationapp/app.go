package notificationapp

import (
	"context"
	"fmt"
	"log/slog"
	http2 "net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/pubsub/redispubsub"
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
	cfg              Config
	trap             <-chan os.Signal
	logger           *slog.Logger
	clientHTTPServer http.ClientServer
	adminHTTPServer  http.AdminServer
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal, re *redis.Adapter, pg *postgres.Database) Application {
	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		panic(tErr)
	}

	cfg.Notification.PingPeriod = (cfg.Notification.PongWait * 9) / 10

	cache := cachemanager.New(re)

	notificationVld := service.NewValidate(t)
	notificationRepo := postgres2.New(pg)

	pubSub := redispubsub.New(re)

	hub := service.NewHub(cfg.Notification, logger, pubSub)
	upgrader := websocket.NewGorillaUpgrader(cfg.Websocket, checkOrigin(cfg.ClientHTTPServer.Cors.AllowOrigins, logger))
	notificationSvc := service.New(cfg.Notification, notificationVld, cache, notificationRepo, logger, hub, pubSub)

	handler := http.NewHandler(notificationSvc, t, upgrader)
	clientHTTPServer := http.NewClientServer(httpserver.New(cfg.ClientHTTPServer, logger), handler, cfg.GetUserIDURL, logger)

	adminHTTPServer := http.NewAdminServer(httpserver.New(cfg.AdminHTTPServer, logger), handler, logger)

	return Application{
		cfg:              cfg,
		trap:             trap,
		logger:           logger,
		clientHTTPServer: clientHTTPServer,
		adminHTTPServer:  adminHTTPServer,
	}
}

func (a Application) Start() {
	clientHTTPServerChan := make(chan error, 1)
	adminHTTPServerChan := make(chan error, 1)

	go func() {
		a.logger.Info(fmt.Sprintf("client http server started on %d port", a.cfg.ClientHTTPServer.Port))

		if sErr := a.clientHTTPServer.Start(); sErr != nil {
			clientHTTPServerChan <- sErr
		}
	}()

	go func() {
		a.logger.Info(fmt.Sprintf("admin http server started on %d", a.cfg.AdminHTTPServer.Port))

		if sErr := a.adminHTTPServer.Start(); sErr != nil {
			adminHTTPServerChan <- sErr
		}
	}()

	select {
	case err := <-clientHTTPServerChan:
		a.logger.Error(fmt.Sprintf("error in client http server on %d", a.cfg.ClientHTTPServer.Port), slog.String("error", err.Error()))
	case err := <-adminHTTPServerChan:
		a.logger.Error(fmt.Sprintf("error in admin http server on %d", a.cfg.AdminHTTPServer.Port), slog.String("error", err.Error()))
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

	if sErr := a.clientHTTPServer.Stop(ctx); sErr != nil {
		a.logger.Error("client http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}

	if sErr := a.adminHTTPServer.Stop(ctx); sErr != nil {
		a.logger.Error("admin http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}
}

func checkOrigin(allowedOrigins []string, logger *slog.Logger) func(r *http2.Request) bool {
	return func(r *http2.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			logger.Warn("ws connection attempt without header")

			// TODO: change it to false in production
			return true
		}

		if len(allowedOrigins) == 0 {
			logger.Debug("accepting all origins because allowedOrigins is empty", slog.String("origin", origin))

			return true
		}

		u, pErr := url.Parse(origin)
		if pErr != nil {
			logger.Warn("invalid origin header", slog.String("origin", origin))

			return false
		}

		hostname := u.Hostname()
		for _, allowed := range allowedOrigins {
			if strings.HasPrefix(allowed, "*.") {
				domain := strings.TrimPrefix(allowed, "*.")
				if strings.HasSuffix(hostname, domain) {
					logger.Debug("origin accepted (wildcard match)",
						slog.String("origin", origin),
						slog.String("pattern", allowed))

					return true
				}
			} else if hostname == allowed || origin == allowed {
				logger.Debug("origin accepted (exact match)",
					slog.String("origin", origin),
					slog.String("allowed", allowed))

				return true
			}
		}

		logger.Warn("origin rejected", slog.String("origin", origin))

		return false
	}
}
