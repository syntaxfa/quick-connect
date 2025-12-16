package http

import (
	"context"
	"log/slog"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/storageapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type Server struct {
	httpServer httpserver.Server
	handler    Handler
	logger     *slog.Logger
}

func New(httpServer httpserver.Server, handler Handler, logger *slog.Logger) Server {
	return Server{
		httpServer: httpServer,
		handler:    handler,
		logger:     logger,
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

	s.httpServer.Router.GET("health-check", s.handler.healthCheck)

	fileGR := s.httpServer.Router.Group("files")
	fileGR.GET("/*", s.handler.ServeFile)
}

func (s Server) registerSwagger() {
	docs.SwaggerInfostorage.Title = "Storage API"
	docs.SwaggerInfostorage.Description = "Storage restfull API documentation"
	docs.SwaggerInfostorage.Version = "1.x.x"

	s.httpServer.Router.GET("/swagger/*any", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("storage")))
}
