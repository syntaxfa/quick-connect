package http

import (
	"context"

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
}

func (s Server) registerSwagger() {
}
