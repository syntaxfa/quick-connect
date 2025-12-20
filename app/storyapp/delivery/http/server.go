package http

import (
	"context"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/storyapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/types"
)

type Server struct {
	httpServer httpserver.Server
	handler    Handler
	authMid    *auth.Middleware
}

func New(httpServer httpserver.Server, handler Handler, authMid *auth.Middleware) Server {
	return Server{
		httpServer: httpServer,
		handler:    handler,
		authMid:    authMid,
	}
}

func (s Server) Start() error {
	s.registerRoutes()

	return s.httpServer.Start()
}

func (s Server) Stop(ctx context.Context) error {
	return s.httpServer.Stop(ctx)
}

func (s Server) registerRoutes() {
	s.registerSwagger()

	s.httpServer.Router.GET("health-check", s.handler.healthCheck)

	storyGr := s.httpServer.Router.Group("stories")
	storyGr.POST("", s.handler.AddStory, s.authMid.RequireAuth,
		s.authMid.RequireRole([]types.Role{types.RoleSuperUser, types.RoleStory}))
}

func (s Server) registerSwagger() {
	docs.SwaggerInfostory.Title = "Story API"
	docs.SwaggerInfostory.Description = "Story restfull API documentation"
	docs.SwaggerInfostory.Version = "1.x.x"

	s.httpServer.Router.GET("/swagger/*any", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("story")))
}
