package http

import (
	"context"

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

	// TODO: Add authentication middleware here
	// protectedGroup := s.httpserver.Router.Group("", AuthMiddleware)
	protectedGroup := s.httpserver.Router.Group("")

	// Dashboard - Main hub
	protectedGroup.GET("/dashboard", s.handler.ShowDashboard)

	// Service routes - these load content via HTMX
	protectedGroup.GET("/support", s.handler.ShowSupportService)
	protectedGroup.GET("/notification", s.handler.ShowNotificationService)
	protectedGroup.GET("/story", s.handler.ShowStoryService)

	// Users management routes
	protectedGroup.GET("/users", s.handler.ShowUsers)
	protectedGroup.GET("/users/search", s.handler.SearchUsers)
	protectedGroup.GET("/users/export", s.handler.ExportUsers)
	protectedGroup.GET("/users/create", s.handler.ShowCreateUserForm)
}

func (s Server) registerSwagger() {
}
