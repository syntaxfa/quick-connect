package storageapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/storage/aws"
	"github.com/syntaxfa/quick-connect/adapter/storage/local"
	grpcdelivery "github.com/syntaxfa/quick-connect/app/storageapp/delivery/grpc"
	"github.com/syntaxfa/quick-connect/app/storageapp/delivery/http"
	postgres2 "github.com/syntaxfa/quick-connect/app/storageapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/rolemanager"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type Application struct {
	cfg                Config
	httpServer         http.Server
	logger             *slog.Logger
	trap               <-chan os.Signal
	internalGRPCServer grpcdelivery.InternalServer
}

type PublicKeyService interface {
	GetPublicKey(ctx context.Context, req *empty.Empty, opts ...grpc.CallOption) (*authpb.GetPublicKeyResponse, error)
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal, t *translation.Translate, psqAdapter *postgres.Database,
	publicKeyInternalAd PublicKeyService) (Application, service.Service) {
	const op = "Setup"

	cfg.Service.Driver = cfg.Storage.Driver
	cfg.Service.Bucket = cfg.Storage.AWS.BucketName

	var storage service.Storage
	var storageErr error

	if cfg.Storage.Driver == service.DriverS3 {
		ctx := context.Background()

		storage, storageErr = aws.New(ctx, cfg.Storage.AWS)
		if storageErr != nil {
			logger.Error("can't connect to s3", slog.String("error", storageErr.Error()))

			panic(storageErr)
		}
	} else {
		storage, storageErr = local.New(cfg.Storage.Local, logger)
		if storageErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(storageErr).WithKind(richerror.KindUnexpected), logger)

			panic(storageErr)
		}
	}

	repo := postgres2.New(psqAdapter)
	svc := service.New(cfg.Service, storage, repo, logger)

	var publicKeyAd PublicKeyService
	var managerGRPCClient *grpcclient.Client

	if publicKeyInternalAd != nil {
		publicKeyAd = publicKeyInternalAd
	} else {
		var grpcErr error
		managerGRPCClient, grpcErr = grpcclient.New(cfg.ManagerAppGRPC)
		if grpcErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(grpcErr).WithKind(richerror.KindUnexpected), logger)

			panic(grpcErr)
		}

		publicKeyAd = manager.NewAuthAdapter(managerGRPCClient.Conn())
	}

	resp, pubErr := publicKeyAd.GetPublicKey(context.Background(), nil)
	if pubErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(pubErr).WithKind(richerror.KindUnexpected), logger)

		panic(pubErr)
	}

	jwtValidator := jwtvalidator.New(resp.GetPublicKey(), logger)
	authMid := auth.New(jwtValidator)

	handler := http.NewHandler(svc, t, cfg.Storage.Local.RootPath, cfg.Service.MaxFileSize, logger)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), handler, logger, authMid)

	internalRoleManager := SetupInternalRoleManager()
	authInterceptor := grpcauth.NewAuthInterceptor(jwtValidator, internalRoleManager)
	internalGRPCHandler := grpcdelivery.NewInternalHandler(svc, t, logger)
	internalGRPCServer := grpcdelivery.New(grpcserver.New(cfg.InternalGRPCServer, logger, grpc.UnaryInterceptor(authInterceptor)),
		internalGRPCHandler)

	return Application{
		cfg:                cfg,
		httpServer:         httpServer,
		logger:             logger,
		trap:               trap,
		internalGRPCServer: internalGRPCServer,
	}, svc
}

func (a Application) Start() {
	httpServerChan := make(chan error, 1)
	grpcServerChan := make(chan error, 1)

	go func() {
		a.logger.Info(fmt.Sprintf("http server started on %d", a.cfg.HTTPServer.Port))

		if sErr := a.httpServer.Start(); sErr != nil {
			httpServerChan <- sErr
		}
	}()

	go func() {
		if sErr := a.internalGRPCServer.Start(context.Background()); sErr != nil {
			grpcServerChan <- sErr
		}
	}()

	select {
	case err := <-httpServerChan:
		a.logger.Error(fmt.Sprintf("error in http server on %d", a.cfg.HTTPServer.Port), slog.String("error", err.Error()))
	case err := <-grpcServerChan:
		a.logger.Error(fmt.Sprintf("error in grpc server on %d", a.cfg.InternalGRPCServer.Port), slog.String("error", err.Error()))
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

	a.logger.Info("storage handler app stopped")
}

func (a Application) Stop(ctx context.Context) bool {
	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup

		shutdownWg.Add(1)
		go a.stopHTTPServer(ctx, &shutdownWg)

		shutdownWg.Add(1)
		go a.stopInternalGRPCServer(&shutdownWg)

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

func (a Application) stopInternalGRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()
	a.internalGRPCServer.Stop()
}

func SetupInternalRoleManager() *rolemanager.RoleManager {
	methodRoles := map[string][]types.Role{
		"/storage.StorageInternalService/GetLink":     {types.RoleService},
		"/storage.StorageInternalService/GetFileInfo": {types.RoleService},
	}

	return rolemanager.NewRoleManager(methodRoles)
}
