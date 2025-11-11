package http

import (
	"context"
	"log/slog"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/managerapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/ratelimit"
	"github.com/syntaxfa/quick-connect/types"
)

type Server struct {
	cfg        Config
	httpserver httpserver.Server
	handler    Handler
	authMid    *auth.Middleware
	cache      *cachemanager.CacheManager
	logger     *slog.Logger
}

func New(cfg Config, httpServer httpserver.Server, handler Handler, authMid *auth.Middleware,
	cache *cachemanager.CacheManager, logger *slog.Logger) Server {
	return Server{
		cfg:        cfg,
		httpserver: httpServer,
		handler:    handler,
		authMid:    authMid,
		cache:      cache,
		logger:     logger,
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
	user.PUT("", s.handler.UserUpdate, s.authMid.RequireAuth)
	user.GET("/:userID", s.handler.UserDetail, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleSuperUser}))
	user.DELETE("/:userID", s.handler.UserDelete, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleSuperUser}))
	user.PUT("/:userID", s.handler.UserUpdateFormSuperuser, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleSuperUser}))
	user.POST("/list", s.handler.UserList, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleSuperUser}))
	user.POST("/login", s.handler.UserLogin)
	user.GET("/profile", s.handler.UserProfile, s.authMid.RequireAuth)
	user.POST("/change-password", s.handler.ChangePassword, s.authMid.RequireAuth)

	guest := user.Group("/guest")
	guest.POST("/register", s.handler.RegisterGuestUser,
		ratelimit.ByIPAddressMiddleware(s.cache, s.cfg.RegisterGuestMaxHint, s.cfg.RegisterGuestDurationLimit, s.logger))
	guest.PUT("/update", s.handler.UpdateGuestUser, s.authMid.RequireAuth, s.authMid.RequireRole([]types.Role{types.RoleGuest}),
		ratelimit.ByIPAddressMiddleware(s.cache, s.cfg.UpdateGuestMaxHint, s.cfg.UpdateGuestDurationLimit, s.logger))
}

func (s Server) registerSwagger() {
	docs.SwaggerInfo.Title = "Manager API"
	docs.SwaggerInfo.Description = "Manager restfull API documentation"
	docs.SwaggerInfo.Version = "1.0.0"

	s.httpserver.Router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
