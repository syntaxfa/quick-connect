package http

import (
	"context"

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

	s.httpServer.Router.GET("health-check", s.handler.healthCheck)
}

func (s Server) registerSwagger() {}
