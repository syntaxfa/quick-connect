package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice1/docs"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type Server struct {
	serviceName string
	httpServer  httpserver.Server
	handler     Handler
}

func New(server httpserver.Server, handler Handler, serviceName string) Server {
	return Server{
		httpServer:  server,
		handler:     handler,
		serviceName: serviceName,
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

	s.httpServer.Router.Use(otelecho.Middleware(s.serviceName, otelecho.WithSkipper(func(c echo.Context) bool {
		return c.Path() == "/auth"
	})))

	s.httpServer.Router.GET("/health-check", s.handler.healthCheck)
}

func (s Server) registerSwagger() {
	docs.SwaggerInfo.Title = s.serviceName
	docs.SwaggerInfo.Description = fmt.Sprintf("%s restfull API documentation", s.serviceName)
	docs.SwaggerInfo.Version = "1.0.0"

	s.httpServer.Router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
