package storyapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/storage"
	"github.com/syntaxfa/quick-connect/app/storyapp/delivery/http"
	postgres2 "github.com/syntaxfa/quick-connect/app/storyapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/storyapp/service"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/tokenmanager"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type Application struct {
	cfg               Config
	httpServer        http.Server
	logger            *slog.Logger
	trap              <-chan os.Signal
	storageGRPCClient *grpcclient.Client
	managerGRPCClient *grpcclient.Client
}

type AuthService interface {
	GetPublicKey(ctx context.Context, req *empty.Empty, opts ...grpc.CallOption) (*authpb.GetPublicKeyResponse, error)
	Login(ctx context.Context, req *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.LoginResponse, error)
	TokenRefresh(ctx context.Context, req *authpb.TokenRefreshRequest, opts ...grpc.CallOption) (*authpb.TokenRefreshResponse, error)
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal, t *translation.Translate, psqAd *postgres.Database,
	storageInternalAd service.StorageService, authInternalAd AuthService) (Application, service.Service) {
	const op = "Setup"

	var storageAd service.StorageService
	var storageGRPCClient *grpcclient.Client

	if storageInternalAd != nil {
		storageAd = storageInternalAd
	} else {
		var grpcErr error
		storageGRPCClient, grpcErr = grpcclient.New(cfg.StorageAppGRPC, grpc.WithUnaryInterceptor(grpcauth.AuthClientInterceptor))
		if grpcErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(grpcErr).WithKind(richerror.KindUnexpected), logger)

			panic(grpcErr)
		}

		storageAd = storage.NewInternalAdapter(storageGRPCClient.Conn())
	}

	var authAd AuthService
	var managerGRPCClient *grpcclient.Client

	if authInternalAd != nil {
		authAd = authInternalAd
	} else {
		var grpcErr error
		managerGRPCClient, grpcErr = grpcclient.New(cfg.ManagerAppGRPC)
		if grpcErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(grpcErr).WithKind(richerror.KindUnexpected), logger)

			panic(grpcErr)
		}

		authAd = manager.NewAuthAdapter(managerGRPCClient.Conn())
	}

	tokenManager := tokenmanager.NewTokenManager(cfg.ServiceAuthInfo.Username, cfg.ServiceAuthInfo.Password, authAd)

	repo := postgres2.New(psqAd)
	vld := service.NewValidate(t)
	svc := service.New(repo, vld, storageAd, tokenManager, logger)

	resp, pubErr := authAd.GetPublicKey(context.Background(), nil)
	if pubErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(pubErr).WithKind(richerror.KindUnexpected), logger)

		panic(pubErr)
	}

	jwtValidator := jwtvalidator.New(resp.GetPublicKey(), logger)
	authMid := auth.New(jwtValidator)

	handler := http.NewHandler(svc, t, logger)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), handler, authMid)

	return Application{
		cfg:               cfg,
		httpServer:        httpServer,
		logger:            logger,
		trap:              trap,
		storageGRPCClient: storageGRPCClient,
		managerGRPCClient: managerGRPCClient,
	}, svc
}

func (a Application) Start() {
	httpServerChan := make(chan error, 1)

	go func() {
		a.logger.Info(fmt.Sprintf("http server started on %d", a.cfg.HTTPServer.Port))

		if sErr := a.httpServer.Start(); sErr != nil {
			httpServerChan <- sErr
		}
	}()

	select {
	case err := <-httpServerChan:
		a.logger.Error(fmt.Sprintf("error in http server on %d", a.cfg.HTTPServer.Port), slog.String("error", err.Error()))
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

	a.logger.Info("story app stopped")
}

func (a Application) Stop(ctx context.Context) bool {
	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup

		shutdownWg.Add(1)
		go a.stopHTTPServer(ctx, &shutdownWg)

		shutdownWg.Add(1)
		go a.stopManagerGRPCClient(ctx, &shutdownWg)

		shutdownWg.Add(1)
		go a.stopStorageGRPCClient(ctx, &shutdownWg)

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

func (a Application) stopHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if sErr := a.httpServer.Stop(ctx); sErr != nil {
		a.logger.ErrorContext(ctx, "http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}
}

func (a Application) stopManagerGRPCClient(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	if a.managerGRPCClient == nil {
		return
	}

	if cErr := a.managerGRPCClient.Close(); cErr != nil {
		a.logger.ErrorContext(ctx, "grpc manager client gracefully shutdown failed", slog.String("error", cErr.Error()))
	}
}

func (a Application) stopStorageGRPCClient(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	if a.storageGRPCClient == nil {
		return
	}

	if cErr := a.storageGRPCClient.Close(); cErr != nil {
		a.logger.ErrorContext(ctx, "grpc storage client gracefully shutdown failed", slog.String("error", cErr.Error()))
	}
}
