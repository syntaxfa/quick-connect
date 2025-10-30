package http

import (
	"context"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/managerapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/types"
)

type Server struct {
	httpserver httpserver.Server
	handler    Handler
	authMid    *auth.Middleware
}

func New(httpServer httpserver.Server, handler Handler, authMid *auth.Middleware) Server {
	return Server{
		httpserver: httpServer,
		handler:    handler,
		authMid:    authMid,
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

	s.httpserver.Router.GET("/health-check", s.handler.healthCheck)

	token := s.httpserver.Router.Group("/tokens")
	token.POST("/refresh", s.handler.RefreshToken)
	token.POST("/validate", s.handler.ValidateToken)

	user := s.httpserver.Router.Group("/users")
	user.POST("", s.handler.CreateUser, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleSuperUser}))
	user.DELETE("/:userID", s.handler.UserDelete, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleSuperUser}))
	user.PUT("/:userID", s.handler.UserUpdateFormSuperuser, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleSuperUser}))
	user.POST("/list", s.handler.UserList, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleSuperUser}))
	user.POST("/login", s.handler.UserLogin)
	user.GET("/profile", s.handler.UserProfile, s.authMid.RequireAuth)
}

func (s Server) registerSwagger() {
	docs.SwaggerInfo.Title = "Manager API"
	docs.SwaggerInfo.Description = "Manager restfull API documentation"
	docs.SwaggerInfo.Version = "1.0.0"

	s.httpserver.Router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
