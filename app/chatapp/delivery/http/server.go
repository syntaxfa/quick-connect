package http

import (
	"context"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/chatapp/docs"
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
	s.registerSwagger()

	s.httpServer.Router.GET("/health-check", s.handler.healthCheck)

	v1 := s.httpServer.Router.Group("/v1")

	chats := v1.Group("/chats")
	chats.GET("/clients", s.handler.WSClientHandler)
	chats.GET("/supports", s.handler.WSSupportHandler)
}

func (s Server) registerSwagger() {
	docs.SwaggerInfo.Title = "CHAT API"
	docs.SwaggerInfo.Description = "Chat restfull API documentation"
	docs.SwaggerInfo.Version = "1.0.0"

	s.httpServer.Router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
