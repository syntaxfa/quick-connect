package chatapp

import (
	"context"
	"fmt"
	"log/slog"
	http2 "net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/syntaxfa/quick-connect/app/chatapp/delivery/http"
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/websocket"
)

const (
	pingPeriodNumerator   = 9
	pingPeriodDenominator = 10
)

type Application struct {
	cfg         Config
	trap        <-chan os.Signal
	chatHandler http.Handler
	logger      *slog.Logger
	httpServer  http.Server
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal) Application {
	cfg.ChatService.PingPeriod = (cfg.ChatService.PongWait * pingPeriodNumerator) / pingPeriodDenominator

	fmt.Printf("%+v", cfg)

	upgrader := websocket.NewGorillaUpgrader(cfg.Websocket, checkOrigin(cfg.HTTPServer.Cors.AllowOrigins, logger))

	chatSvc := service.New(cfg.ChatService, nil, logger)
	chatHandler := http.NewHandler(upgrader, logger, chatSvc)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), chatHandler)

	return Application{
		cfg:         cfg,
		chatHandler: chatHandler,
		logger:      logger,
		httpServer:  httpServer,
		trap:        trap,
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

	a.logger.Info("chat app stopped")
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
		a.logger.ErrorContext(ctx, "http server gracefully shutdown failed", slog.String("error", sErr.Error()))
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
