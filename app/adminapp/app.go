package adminapp

import (
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"google.golang.org/grpc"
	"log/slog"
	"os"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/app/adminapp/delivery/http"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Application struct {
	cfg               Config
	trap              <-chan os.Signal
	logger            *slog.Logger
	httpServer        http.Server
	managerGRPCClient *grpcclient.Client
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal) Application {
	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		logger.Error("failed create new instance of translation", slog.String("error", tErr.Error()))

		panic(tErr)
	}

	managerGRPCClient, grpcErr := grpcclient.New(cfg.ManagerAppGRPC, grpc.WithUnaryInterceptor(grpcauth.AuthClientInterceptor))
	if grpcErr != nil {
		logger.Error("failed to create manager gRPC client", slog.String("error", grpcErr.Error()))

		panic(grpcErr)
	}

	authAdapter := manager.NewAuthAdapter(managerGRPCClient.Conn())

	handler := http.NewHandler(logger, t, authAdapter)

	return Application{
		cfg:               cfg,
		trap:              trap,
		logger:            logger,
		httpServer:        http.New(httpserver.New(cfg.HTTPServer, logger), handler, cfg.TemplatePath),
		managerGRPCClient: managerGRPCClient,
	}
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
		a.logger.Error(fmt.Sprintf("error in http server on port %d", a.cfg.HTTPServer.Port), slog.String("error", err.Error()))
	case <-a.trap:
		a.logger.Info("received shutdown signal!!!")
	}

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	if a.Stop(shutdownTimeoutCtx) {
		a.logger.Info("servers shutdown gracefully")
	} else {
		a.logger.Warn("shutdown timed out, existing application")
	}
}

func (a Application) Stop(ctx context.Context) bool {
	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup
		shutdownWg.Add(1)
		go a.stopHTTPServer(ctx, &shutdownWg)

		shutdownWg.Add(1)
		go a.closeManagerGRPCClient(&shutdownWg)

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
		a.logger.Error("http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}
}

func (a Application) closeManagerGRPCClient(wg *sync.WaitGroup) {
	defer wg.Done()
	if cErr := a.managerGRPCClient.Close(); cErr != nil {
		a.logger.Error("failed to close manager gRPC client", slog.String("error", cErr.Error()))
	}
}
