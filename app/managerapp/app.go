package managerapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	grpcdelivery "github.com/syntaxfa/quick-connect/app/managerapp/delivery/grpc"
	"github.com/syntaxfa/quick-connect/app/managerapp/delivery/http"
	postgres2 "github.com/syntaxfa/quick-connect/app/managerapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/pkg/rolemanager"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
)

type Application struct {
	cfg                Config
	trap               <-chan os.Signal
	logger             *slog.Logger
	httpServer         http.Server
	internalHTTPServer http.InternalServer
	grpcServer         grpcdelivery.Server
	grpcServerInternal grpcdelivery.ServerInternal
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal, psqAdapter *postgres.Database, re *redis.Adapter) (
	Application, tokenservice.Service, userservice.Service, *jwtvalidator.Validator) {
	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		panic(tErr)
	}

	tokenSvc := tokenservice.New(cfg.Token, logger)
	vldUser := userservice.NewValidate(t)
	cache := cachemanager.New(re, logger)

	userRepo := postgres2.New(psqAdapter)
	userSvc := userservice.New(cfg.User, tokenSvc, vldUser, userRepo, userRepo, logger, cache)
	handler := http.NewHandler(t, tokenSvc, userSvc)

	jwtValidator := jwtvalidator.New(cfg.Token.PublicKeyString, logger)
	authMid := auth.New(jwtValidator)
	httpServer := http.New(cfg.Delivery, httpserver.New(cfg.HTTPServer, logger), handler, authMid, cache, logger)

	internalHandler := http.NewInternalHandler(userSvc, t)
	internalHTTPServer := http.NewInternalServer(httpserver.New(cfg.InternalHTTPServer, logger), internalHandler, logger, cfg.APIKey)

	roleManager := SetupRoleManager()
	authInterceptor := grpcauth.NewAuthInterceptor(jwtValidator, roleManager)
	grpcHandler := grpcdelivery.NewHandler(logger, tokenSvc, userSvc, t)
	grpcServer := grpcdelivery.New(grpcserver.New(cfg.GRPCServer, logger, grpc.UnaryInterceptor(authInterceptor)), grpcHandler, logger)

	internalRoleManager := SetupInternalRoleManager()
	internalAuthInterceptor := grpcauth.NewAuthInterceptor(jwtValidator, internalRoleManager)
	grpcHandlerInternal := grpcdelivery.NewHandlerInternal(logger, userSvc, t)
	grpcServerInternal := grpcdelivery.NewServerInternal(grpcserver.New(cfg.GRPCServerInternal, logger,
		grpc.UnaryInterceptor(internalAuthInterceptor)), grpcHandlerInternal, logger)

	return Application{
		cfg:                cfg,
		trap:               trap,
		logger:             logger,
		httpServer:         httpServer,
		grpcServer:         grpcServer,
		internalHTTPServer: internalHTTPServer,
		grpcServerInternal: grpcServerInternal,
	}, tokenSvc, userSvc, jwtValidator
}

func (a Application) Start() {
	httpServerChan := make(chan error, 1)
	internalHTTPServerChan := make(chan error, 1)
	grpcServerChan := make(chan error, 1)
	grpcServerInternalChan := make(chan error, 1)

	go func() {
		a.logger.Info(fmt.Sprintf("http server started on %d", a.cfg.HTTPServer.Port))

		if sErr := a.httpServer.Start(); sErr != nil {
			httpServerChan <- sErr
		}
	}()

	go func() {
		a.logger.Info(fmt.Sprintf("internal http server started on %d", a.cfg.InternalHTTPServer.Port))

		if sErr := a.internalHTTPServer.Start(); sErr != nil {
			internalHTTPServerChan <- sErr
		}
	}()

	go func() {
		if sErr := a.grpcServer.Start(); sErr != nil {
			grpcServerChan <- sErr
		}
	}()

	go func() {
		if sErr := a.grpcServerInternal.Start(); sErr != nil {
			grpcServerInternalChan <- sErr
		}
	}()

	select {
	case err := <-httpServerChan:
		a.logger.Error(fmt.Sprintf("error in http server on %d", a.cfg.HTTPServer.Port), slog.String("error", err.Error()))
	case err := <-internalHTTPServerChan:
		a.logger.Error(fmt.Sprintf("error in internal http server on %d", a.cfg.InternalHTTPServer.Port),
			slog.String("error", err.Error()))
	case err := <-grpcServerChan:
		a.logger.Error(fmt.Sprintf("error in grpc server on %d", a.cfg.GRPCServer.Port), slog.String("error", err.Error()))
	case err := <-grpcServerInternalChan:
		a.logger.Error(fmt.Sprintf("error in groc server internal on %d", a.cfg.GRPCServerInternal.Port),
			slog.String("error", err.Error()))
	case <-a.trap:
		a.logger.Info("received http server shutdown signal!!!")
	}

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	if a.Stop(shutdownTimeoutCtx) {
		a.logger.Info("servers shutdown gracefully")
	} else {
		a.logger.Warn("shutdown timed out, existing application")
	}

	a.logger.Info("manager app stopped")
}

func (a Application) Stop(ctx context.Context) bool {
	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup

		shutdownWg.Add(1)
		go a.StopHTTPServer(ctx, &shutdownWg)

		shutdownWg.Add(1)
		go a.StopInternalHTTPServer(ctx, &shutdownWg)

		shutdownWg.Add(1)
		go a.StopGRPCServer(&shutdownWg)

		shutdownWg.Add(1)
		go a.StopGRCPServerInternal(&shutdownWg)

		shutdownWg.Wait()
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
		return true
	case <-ctx.Done():
		return false
	}
}

func (a Application) StopHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if sErr := a.httpServer.Stop(ctx); sErr != nil {
		a.logger.ErrorContext(ctx, "http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}
}

func (a Application) StopInternalHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if sErr := a.internalHTTPServer.Stop(ctx); sErr != nil {
		a.logger.ErrorContext(ctx, "internal http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}
}

func (a Application) StopGRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()
	a.grpcServer.Stop()
}

func (a Application) StopGRCPServerInternal(wg *sync.WaitGroup) {
	defer wg.Done()
	a.grpcServerInternal.Stop()
}

func SetupRoleManager() *rolemanager.RoleManager {
	methodRoles := map[string][]types.Role{
		// AuthService
		"/manager.AuthService/GetPublicKey": {},
		"/manager.AuthService/Login":        {},
		"/manager.AuthService/TokenRefresh": {},
		"/manager.AuthService/TokenVerify":  {},
		// UserService
		"/manager.UserService/CreateUser":              {types.RoleSuperUser},
		"/manager.UserService/UserDelete":              {types.RoleSuperUser},
		"/manager.UserService/UserDetail":              {types.RoleSuperUser},
		"/manager.UserService/UserUpdateFromSuperuser": {types.RoleSuperUser},
		"/manager.UserService/UserList":                {types.RoleSuperUser},
		"/manager.UserService/UserProfile":             types.AdminRoles,
		"/manager.UserService/UserUpdateFromOwn":       types.AdminRoles,
		"/manager.UserService/UserInfo":                types.AllUserRole,
		"/manager.UserService/UserChangePassword":      types.AdminRoles,
	}

	return rolemanager.NewRoleManager(methodRoles)
}

func SetupInternalRoleManager() *rolemanager.RoleManager {
	methodRoles := map[string][]types.Role{
		// UserService
		"/manager.UserInternalService/UserInfo": {types.RoleService},
	}

	return rolemanager.NewRoleManager(methodRoles)
}
