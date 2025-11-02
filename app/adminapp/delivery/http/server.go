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

	// --- Users Management Routes ---
	// This replaces your old route structure
	userGr := rootGr.Group("/users")
	{
		// GET /users
		// Renders the main page shell (users_page.html)
		userGr.GET("", s.handler.ShowUsersPage)

		// GET /users/list
		// This is the new HTMX partial route for searching, pagination, and initial load.
		// It replaces the old /search route.
		userGr.GET("/list", s.handler.ListUsersPartial)

		// POST /users/:id/delete
		// Handles user deletion via HTMX (from users_list_partial.html)
		userGr.POST("/:id/delete", s.handler.DeleteUser)

		// --- Routes for Modals (from templates) ---
		// TODO: Implement these handlers

		// GET /users/create
		// Shows the "Add User" modal
		//userGr.GET("/create", s.handler.ShowCreateUserModal) // TODO: Create h.ShowCreateUserModal

		// GET /users/:id/edit
		// Shows the "Edit User" modal
		//userGr.GET("/:id/edit", h.ShowEditUserModal) // TODO: Create h.ShowEditUserModal

		// GET /users/:id/details
		// Shows the "View Details" modal
		//userGr.GET("/:id/details", h.ShowUserDetailsModal) // TODO: Create h.ShowUserDetailsModal

		// --- Other Actions ---

		// GET /users/export
		// (Handler not yet implemented, but route is defined in template)
		// userGr.GET("/export", h.ExportUsers) // TODO: Create h.ExportUsers handler
	}

	// Profile Group
	profileGr := rootGr.Group("/profile")
	profileGr.GET("", s.handler.ShowProfilePage)
	profileGr.PUT("", s.handler.UpdateProfile)
	profileGr.GET("/view", s.handler.ShowProfileView)
	profileGr.GET("/edit", s.handler.ShowProfileEditForm)

	// Settings
	settingGr := rootGr.Group("/settings")
	settingGr.GET("", s.handler.ShowSettingsPage)
}

func (s Server) registerSwagger() {
}

func grpcContext(c echo.Context) context.Context {
	token := c.Get(types.AuthorizationKey)

	return context.WithValue(c.Request().Context(), types.AuthorizationKey, token)
}
