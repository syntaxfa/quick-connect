package http

import (
	"context"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/notificationapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type Server struct {
	httpServer httpserver.Server
	handler    Handler
}

func New(httpServer httpserver.Server, handler Handler) Server {
	return Server{
		httpServer: httpServer,
		handler:    handler,
	}
}

func (s Server) Start() error {
	s.registerRoutes()

	return s.httpServer.Start()
}

func (s Server) Stop(ctx context.Context) error {
	return s.httpServer.Stop(ctx)
}

func (s Server) registerRoutes() {
	s.registerSwagger()

	s.httpServer.Router.GET("/health-check", s.handler.healthCheck)

	notifications := s.httpServer.Router.Group("/notifications")
	notifications.POST("", s.handler.sendNotification)
	notifications.GET("/ws", s.handler.wsNotification)
}

func (s Server) registerSwagger() {
	docs.SwaggerInfo.Title = "Notification API"
	docs.SwaggerInfo.Description = "Notification restfull API documentation"
	docs.SwaggerInfo.Version = "1.0.0"

	s.httpServer.Router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
