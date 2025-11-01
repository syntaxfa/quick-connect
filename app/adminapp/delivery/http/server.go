package http

import (
	"context"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"

	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type Server struct {
	httpserver   httpserver.Server
	handler      Handler
	jwtValidator *jwtvalidator.Validator
}

func New(httpServer httpserver.Server, handler Handler, templatePath string, jwtValidator *jwtvalidator.Validator) Server {
	renderer := NewTemplateRenderer(templatePath)
	httpServer.Router.Renderer = renderer

	return Server{
		httpserver:   httpServer,
		handler:      handler,
		jwtValidator: jwtValidator,
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

	rootGr := s.httpserver.Router.Group("", setTokenToRequestContextMiddleware(s.jwtValidator, s.handler.authAd, "/login", s.handler.logger))

	// authGroup
	authGroup := rootGr.Group("")
	authGroup.GET("/login", s.handler.ShowLoginPage)
	authGroup.POST("/login", s.handler.Login)
	authGroup.GET("/logout", s.handler.Logout)

	// Dashboard - Main hub
	dashGr := rootGr.Group("")
	dashGr.GET("/dashboard", s.handler.ShowDashboard)
	dashGr.GET("/support", s.handler.ShowSupportService)
	dashGr.GET("/notification", s.handler.ShowNotificationService)
	dashGr.GET("/story", s.handler.ShowStoryService)

	// Users management routes
	userGr := rootGr.Group("/users")
	userGr.GET("", s.handler.ShowUsers)
	userGr.GET("/search", s.handler.SearchUsers)
	userGr.GET("/export", s.handler.ExportUsers)
	userGr.GET("/create", s.handler.ShowCreateUserForm)

	// Profile
	profileGroup := rootGr.Group("/profile")
	profileGroup.GET("", s.handler.ShowProfilePage)
	profileGroup.PUT("", s.handler.UpdateProfile)
	profileGroup.GET("/view", s.handler.ShowProfileView)
	profileGroup.GET("/edit", s.handler.ShowProfileEditForm)
}

func (s Server) registerSwagger() {
}
