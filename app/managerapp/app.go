package managerapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/app/managerapp/delivery/grpc"
	"github.com/syntaxfa/quick-connect/app/managerapp/delivery/http"
	postgres2 "github.com/syntaxfa/quick-connect/app/managerapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Application struct {
	cfg        Config
	trap       <-chan os.Signal
	logger     *slog.Logger
	httpServer http.Server
	grpcServer grpc.Server
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal, psqAdapter *postgres.Database) Application {
	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		panic(tErr)
	}

	tokenSvc := tokenservice.New(cfg.Token, logger)
	vldUser := userservice.NewValidate(t)

	userRepo := postgres2.New(psqAdapter)
	userSvc := userservice.New(tokenSvc, vldUser, userRepo, logger)
	handler := http.NewHandler(t, tokenSvc, userSvc)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), handler)

	grpcHandler := grpc.NewHandler(logger)
	grpcServer := grpc.New(grpcserver.New(cfg.GRPCServer, logger), grpcHandler, logger)

	return Application{
		cfg:        cfg,
		trap:       trap,
		logger:     logger,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
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
		if sErr := a.grpcServer.Start(); sErr != nil {
			grpcServerChan <- sErr
		}
	}()

	select {
	case err := <-httpServerChan:
		a.logger.Error(fmt.Sprintf("error in http server on %d", a.cfg.HTTPServer.Port), slog.String("error", err.Error()))
	case err := <-grpcServerChan:
		a.logger.Error(fmt.Sprintf("error in grpc server on %d", a.cfg.GRPCServer.Port), slog.String("error", err.Error()))
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

		shutdownWg.Wait()
		close(shutdownDone)
	}()

	go func() {
		a.grpcServer.Stop()
		a.logger.Info("grpc server gracefully stop")
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
		a.logger.Error("http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}
}
