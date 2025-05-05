package http

import (
	"context"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type Server struct {
	httpServer httpserver.Server
	handler    Handler
}

func New(server httpserver.Server, handler Handler) Server {
	return Server{
		httpServer: server,
		handler:    handler,
	}
}

func (s Server) Start() error {
	s.RegisterRoutes()

	return s.httpServer.Start()
}

func (s Server) Stop(ctx context.Context) error {
	return s.httpServer.Stop(ctx)
}

func (s Server) RegisterRoutes() {
	s.httpServer.Router.GET("/health-check", s.handler.healthCheck)

	//v1 := s.httpServer.Router.Group("/v1")
}
