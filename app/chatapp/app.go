package chatapp

import (
	"context"
	"fmt"
	"log/slog"
	http2 "net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/syntaxfa/quick-connect/app/chatapp/delivery/http"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type Application struct {
	cfg         Config
	chatHandler http.Handler
	trap        <-chan os.Signal
	logger      *slog.Logger
	httpServer  http.Server
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal) Application {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 0,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      checkOrigin(cfg.HTTPServer.Cors.AllowOrigins, logger),
	}

	chatHandler := http.NewHandler(upgrader, logger)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), chatHandler)

	return Application{
		cfg:         cfg,
		chatHandler: chatHandler,
		trap:        trap,
		logger:      logger,
		httpServer:  httpServer,
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
