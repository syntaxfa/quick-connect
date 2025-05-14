package http

import (
	"context"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/managerapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type Server struct {
	httpserver httpserver.Server
	handler    Handler
}

func New(httpServer httpserver.Server, handler Handler) Server {
	return Server{
		httpserver: httpServer,
		handler:    handler,
	}
}

func (s Server) Start() error {
	s.registerRoutes()

	return s.httpserver.Start()
}

func (s Server) Stop(ctx context.Context) error {
	return s.httpserver.Stop(ctx)
}

func (s Server) registerRoutes() {
	s.registerSwagger()

	s.httpserver.Router.GET("/health-check", s.handler.healthCheck)

	token := s.httpserver.Router.Group("/tokens")
	token.POST("/refresh", s.handler.RefreshToken)
	token.POST("/validate", s.handler.ValidateToken)
}

func (s Server) registerSwagger() {
	docs.SwaggerInfo.Title = "Manager API"
	docs.SwaggerInfo.Description = "Manager restfull API documentation"
	docs.SwaggerInfo.Version = "1.0.0"

	s.httpserver.Router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
