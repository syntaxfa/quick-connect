package adminapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/chat"
	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/app/adminapp/delivery/http"
	"github.com/syntaxfa/quick-connect/app/adminapp/service"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"google.golang.org/grpc"
)

type Application struct {
	cfg               Config
	trap              <-chan os.Signal
	logger            *slog.Logger
	httpServer        http.Server
	managerGRPCClient *grpcclient.Client
	chatGRPCClient    *grpcclient.Client
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal, t *translation.Translate, authLocalAdapter service.AuthService,
	userLocalAdapter service.UserService, conversationLocalAdapter service.ConversationService) Application {
	const op = "Setup"

	var authAdapter service.AuthService
	var userAdapter service.UserService
	var managerGRPCClient *grpcclient.Client

	if authLocalAdapter != nil || userLocalAdapter != nil {
		authAdapter = authLocalAdapter
		userAdapter = userLocalAdapter
	} else {
		var grpcErr error
		managerGRPCClient, grpcErr = grpcclient.New(cfg.ManagerAppGRPC, grpc.WithUnaryInterceptor(grpcauth.AuthClientInterceptor))
		if grpcErr != nil {
			logger.Error("failed to create manager gRPC client", slog.String("error", grpcErr.Error()))

			panic(grpcErr)
		}

		authAdapter = manager.NewAuthAdapter(managerGRPCClient.Conn())
		userAdapter = manager.NewUserAdapter(managerGRPCClient.Conn())
	}

	var conversationAdapter service.ConversationService
	var chatGRPCClient *grpcclient.Client
	if conversationLocalAdapter != nil {
		conversationAdapter = conversationLocalAdapter
	} else {
		var chatErr error
		chatGRPCClient, chatErr = grpcclient.New(cfg.ChatAppGRPC, grpc.WithUnaryInterceptor(grpcauth.AuthClientInterceptor))
		if chatErr != nil {
			logger.Error("failed to create manager gRPC client", slog.String("error", chatErr.Error()))

			panic(chatErr)
		}

		conversationAdapter = chat.NewConversationAdapter(chatGRPCClient.Conn())
	}

	handler := http.NewHandler(logger, t, authAdapter, userAdapter, conversationAdapter)

	getPuResp, gpuErr := authAdapter.GetPublicKey(context.Background(), nil)
	if gpuErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(gpuErr).WithKind(richerror.KindUnexpected), logger)

		panic(gpuErr)
	}

	jwtValidator := jwtvalidator.New(getPuResp.GetPublicKey(), logger)

	return Application{
		cfg:               cfg,
		trap:              trap,
		logger:            logger,
		httpServer:        http.New(httpserver.New(cfg.HTTPServer, logger), handler, cfg.TemplatePath, jwtValidator),
		managerGRPCClient: managerGRPCClient,
		chatGRPCClient:    chatGRPCClient,
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
		go a.stopManagerGRPCClient(&shutdownWg)

		shutdownWg.Add(1)
		go a.stopChatGRPCClient(&shutdownWg)

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

func (a Application) stopManagerGRPCClient(wg *sync.WaitGroup) {
	defer wg.Done()
	if a.managerGRPCClient == nil {
		return
	}

	if cErr := a.managerGRPCClient.Close(); cErr != nil {
		a.logger.Error("failed to close manager gRPC client", slog.String("error", cErr.Error()))
	}
}

func (a Application) stopChatGRPCClient(wg *sync.WaitGroup) {
	defer wg.Done()
	if a.chatGRPCClient == nil {
		return
	}

	if cErr := a.chatGRPCClient.Close(); cErr != nil {
		a.logger.Error("failed to close chat gRPC client", slog.String("error", cErr.Error()))
	}
}
