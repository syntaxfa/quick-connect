package http

import (
	"context"
	"log/slog"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/notificationapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type AdminServer struct {
	httpServer httpserver.Server
	handler    Handler
	logger     *slog.Logger
}

func NewAdminServer(httpServer httpserver.Server, handler Handler, logger *slog.Logger) AdminServer {
	return AdminServer{
		httpServer: httpServer,
		handler:    handler,
		logger:     logger,
	}
}

func (s AdminServer) Start() error {
	s.registerRoutes()

	return s.httpServer.Start()
}

func (s AdminServer) Stop(ctx context.Context) error {
	return s.httpServer.Stop(ctx)
}

func (s AdminServer) registerRoutes() {
	s.registerSwagger()

	s.httpServer.Router.GET("/health-check", s.handler.healthCheck)

	v1 := s.httpServer.Router.Group("/v1")

	notifications := v1.Group("/notifications")

	notifications.POST("", s.handler.sendNotification)
}

func (s AdminServer) registerSwagger() {
	docs.SwaggerInfo.Title = "Notification Admin API"
	docs.SwaggerInfo.Description = "Notification admin restfull API documentation"
	docs.SwaggerInfo.Version = "1.0.0"

	s.httpServer.Router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
