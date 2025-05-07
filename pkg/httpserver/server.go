package httpserver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Cors struct {
	AllowOrigins []string `koanf:"allow_origins"`
}

type Config struct {
	Port int  `koanf:"port"`
	Cors Cors `koang:"cors"`
}

type Server struct {
	Router *echo.Echo
	cfg    Config
	log    *slog.Logger
}

func New(cfg Config, log *slog.Logger) Server {
	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRequestID:     true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			errMsg := ""
			if v.Error != nil {
				errMsg = v.Error.Error()
			}

			log.Info("http-server",
				slog.String("Latency", v.Latency.String()),
				slog.String("Protocol", v.Protocol),
				slog.String("RemoteIP", v.RemoteIP),
				slog.String("Host", v.Host),
				slog.String("Method", v.Method),
				slog.String("URI", v.URI),
				slog.String("URLPath", v.URIPath),
				slog.String("RequestID", v.RequestID),
				slog.String("UserAgent", v.UserAgent),
				slog.Int("Status", v.Status),
				slog.String("ErrorMsg", errMsg),
				slog.String("ContentLength", v.ContentLength),
				slog.Int64("ResponseSize", v.ResponseSize),
			)

			return nil
		},
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.Cors.AllowOrigins,
	}))

	return Server{
		Router: e,
		cfg:    cfg,
		log:    log,
	}
}

func (s Server) Start() error {
	s.log.Info("http server started", slog.Int("port", s.cfg.Port))

	return s.Router.Start(fmt.Sprintf(":%d", s.cfg.Port))
}

func (s Server) Stop(ctx context.Context) error {
	s.log.Info("http server gracefully shutdown", slog.Int("port", s.cfg.Port))

	return s.Router.Shutdown(ctx)
}
