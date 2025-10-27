package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type Server struct {
	httpserver httpserver.Server
	handler    Handler
}

func New(httpServer httpserver.Server, handler Handler, templatePath string) Server {
	renderer := NewTemplateRenderer(templatePath)
	httpServer.Router.Renderer = renderer

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

	s.httpserver.Router.Static("/static", "app/adminapp/static")

	templateRout := s.httpserver.Router.Group("")
	templateRout.GET("/login", s.handler.ShowLoginPage)

	authGroup := s.httpserver.Router.Group("")
	authGroup.POST("/login", s.handler.Login)

	s.httpserver.Router.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "<div id='ping-result'>Pong! HTMX works.</div>")
	})
}

func (s Server) registerSwagger() {
}
