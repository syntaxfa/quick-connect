package http

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/types"
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

	// auth Group
	authGr := rootGr.Group("")
	authGr.GET("/login", s.handler.ShowLoginPage)
	authGr.POST("/login", s.handler.Login)
	authGr.GET("/logout", s.handler.Logout)
	authGr.GET("/logout/confirm", s.handler.ShowLogoutConfirm)

	// Dashboard - Main hub
	dashGr := rootGr.Group("")
	dashGr.GET("/dashboard", s.handler.ShowDashboard)
	dashGr.GET("/support", s.handler.ShowSupportService)
	dashGr.GET("/notification", s.handler.ShowNotificationService)
	dashGr.GET("/story", s.handler.ShowStoryService)

	// Users Management Group
	userGr := rootGr.Group("/users")
	userGr.GET("", s.handler.ShowUsersPage)         // Renders the main page shell (users_page.html)
	userGr.GET("/list", s.handler.ListUsersPartial) // This is the new HTMX partial route for searching, pagination, and initial load.
	userGr.GET("/delete/confirm", s.handler.ShowDeleteUserConfirm)
	userGr.POST("/:id/delete", s.handler.DeleteUser)
	userGr.GET("/:id/detail", s.handler.DetailUser)
	userGr.GET("/:id/edit", s.handler.ShowEditUserModal)
	userGr.POST("/:id/update", s.handler.UpdateUser)
	userGr.GET("/create", s.handler.ShowCreateUserModal)
	userGr.POST("/create", s.handler.CreateUser)

	// Profile Group
	profileGr := rootGr.Group("/profile")
	profileGr.GET("", s.handler.ShowProfilePage)
	profileGr.PUT("", s.handler.UpdateProfile)
	profileGr.GET("/view", s.handler.ShowProfileView)
	profileGr.GET("/edit", s.handler.ShowProfileEditForm)

	// Settings Group
	settingGr := rootGr.Group("/settings")
	settingGr.GET("", s.handler.ShowSettingsPage)
}

func (s Server) registerSwagger() {
}

func grpcContext(c echo.Context) context.Context {
	token := c.Get(string(types.AuthorizationKey))

	return context.WithValue(c.Request().Context(), types.AuthorizationKey, token)
}
