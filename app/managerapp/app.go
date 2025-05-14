package managerapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/app/managerapp/delivery/http"
	postgres2 "github.com/syntaxfa/quick-connect/app/managerapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Application struct {
	cfg        Config
	trap       <-chan os.Signal
	logger     *slog.Logger
	httpServer http.Server
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal) Application {
	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		panic(tErr)
	}

	postgresAdapter := postgres.New(cfg.Postgres)

	tokenSvc := tokenservice.New(cfg.Token, logger)
	vldUser := userservice.NewValidate(t)

	userRepo := postgres2.New(postgresAdapter)
	userSvc := userservice.New(tokenSvc, vldUser, userRepo)
	handler := http.NewHandler(t, tokenSvc, userSvc)
	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), handler)

	return Application{
		cfg:        cfg,
		trap:       trap,
		logger:     logger,
		httpServer: httpServer,
	}
}

func (a Application) Start() {
	go func() {
		a.logger.Info(fmt.Sprintf("http server started in %d", a.cfg.HTTPServer.Port))

		if sErr := a.httpServer.Start(); sErr != nil {
			a.logger.Error(fmt.Sprintf("error in http server on %d", a.cfg.HTTPServer.Port), slog.String("error", sErr.Error()))
		}
		a.logger.Info(fmt.Sprintf("http server stopped %d", a.cfg.HTTPServer.Port))
	}()

	<-a.trap

	a.logger.Info("shutdown signal received")

	ctx, cancelFunc := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancelFunc()

	a.Stop(ctx)
}

func (a Application) Stop(ctx context.Context) {
	if sErr := a.httpServer.Stop(ctx); sErr != nil {
		a.logger.Error("http server gracefully shutdown failed", slog.String("error", sErr.Error()))
	}

	a.logger.Info("http server gracefully shutdown")
}
