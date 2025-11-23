package chatapp

import (
	"context"
	"fmt"
	"log/slog"
	http2 "net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/pubsub/redispubsub"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	grpcdelivery "github.com/syntaxfa/quick-connect/app/chatapp/delivery/grpc"
	"github.com/syntaxfa/quick-connect/app/chatapp/delivery/http"
	postgres2 "github.com/syntaxfa/quick-connect/app/chatapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
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
	"github.com/syntaxfa/quick-connect/pkg/websocket"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
)

const (
	pingPeriodNumerator   = 9
	pingPeriodDenominator = 10
)

type Application struct {
	cfg               Config
	trap              <-chan os.Signal
	chatHandler       http.Handler
	logger            *slog.Logger
	httpServer        http.Server
	managerGRPCClient *grpcclient.Client
	grpcServer        grpcdelivery.Server
	chatHub           *service.Hub       // To control the hub lifecycle
	mainCtx           context.Context    // Main context for background services
	mainCancel        context.CancelFunc // Function to cancel main context
}

func Setup(cfg Config, logger *slog.Logger, trap <-chan os.Signal, psqAdapter *postgres.Database, re *redis.Adapter) Application {
	const op = "Setup"

	mainCtx, mainCancel := context.WithCancel(context.Background())

	cfg.ChatService.PingPeriod = (cfg.ChatService.PongWait * pingPeriodNumerator) / pingPeriodDenominator

	upgrader := websocket.NewGorillaUpgrader(cfg.Websocket, checkOrigin(cfg.HTTPServer.Cors.AllowOrigins, logger))

	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(tErr).WithKind(richerror.KindUnexpected), logger)

		panic(tErr)
	}

	pubsubClient := redispubsub.New(re)

	chatRepo := postgres2.New(psqAdapter)
	vld := service.NewValidate(t)

	chatHub := service.NewHub(cfg.ChatService, logger, pubsubClient)

	chatSvc := service.New(cfg.ChatService, chatRepo, chatHub, pubsubClient, logger, vld)
	chatHandler := http.NewHandler(mainCtx, upgrader, logger, chatSvc, t)

	managerGRPCClient, grpcErr := grpcclient.New(cfg.ManagerAppGRPC)
	if grpcErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(grpcErr).WithKind(richerror.KindUnexpected), logger)

		panic(grpcErr)
	}

	authAd := manager.NewAuthAdapter(managerGRPCClient.Conn())
	resp, pubErr := authAd.GetPublicKey(context.Background())
	if pubErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(pubErr).WithKind(richerror.KindUnexpected), logger)

		panic(pubErr)
	}

	jwtValidator := jwtvalidator.New(resp.GetPublicKey(), logger)
	authMid := auth.New(jwtValidator)

	httpServer := http.New(httpserver.New(cfg.HTTPServer, logger), chatHandler, authMid)

	roleManager := setupRoleManager()
	authInterceptor := grpcauth.NewAuthInterceptor(jwtValidator, roleManager)
	grpcHandler := grpcdelivery.NewHandler(chatSvc, t, logger)
	grpcServer := grpcdelivery.New(grpcserver.New(cfg.GRPCServer, logger, grpc.UnaryInterceptor(authInterceptor)), grpcHandler, logger)

	return Application{
		cfg:               cfg,
		chatHandler:       chatHandler,
		logger:            logger,
		httpServer:        httpServer,
		trap:              trap,
		managerGRPCClient: managerGRPCClient,
		grpcServer:        grpcServer,
		chatHub:           chatHub,
		mainCtx:           mainCtx,
		mainCancel:        mainCancel,
	}
}

func (a Application) Start() {
	httpServerChan := make(chan error, 1)
	grpcServerChan := make(chan error, 1)

	go a.chatHub.Run(a.mainCtx)
	a.logger.Info("chat hub started")

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

	a.logger.Info("chat app stopped")
}

func (a Application) Stop(ctx context.Context) bool {
	// This cancels the mainCtx passed to hub.Run()
	a.mainCancel()

	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup

		shutdownWg.Add(1)
		go a.StopHTTPServer(ctx, &shutdownWg)

		shutdownWg.Add(1)
		go a.StopGRPCClient(ctx, &shutdownWg)

		shutdownWg.Add(1)
		go a.StopGRPCServer(&shutdownWg)

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

func (a Application) StopGRPCClient(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if cErr := a.managerGRPCClient.Close(); cErr != nil {
		a.logger.ErrorContext(ctx, "http server gracefully shutdown failed", slog.String("error", cErr.Error()))
	}
}

func (a Application) StopGRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()
	a.grpcServer.Stop()
}

func checkOrigin(allowedOrigins []string, logger *slog.Logger) func(r *http2.Request) bool {
	return func(r *http2.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			logger.Warn("ws connection attempt without header")

			// TODO: change it to false in production
			return true
		}

		if len(allowedOrigins) == 0 {
			logger.Debug("accepting all origins because allowedOrigins is empty", slog.String("origin", origin))

			return true
		}

		u, pErr := url.Parse(origin)
		if pErr != nil {
			logger.Warn("invalid origin header", slog.String("origin", origin))

			return false
		}

		hostname := u.Hostname()
		for _, allowed := range allowedOrigins {
			if strings.HasPrefix(allowed, "*.") {
				domain := strings.TrimPrefix(allowed, "*.")
				if strings.HasSuffix(hostname, domain) {
					logger.Debug("origin accepted (wildcard match)",
						slog.String("origin", origin),
						slog.String("pattern", allowed))

					return true
				}
			} else if hostname == allowed || origin == allowed {
				logger.Debug("origin accepted (exact match)",
					slog.String("origin", origin),
					slog.String("allowed", allowed))

				return true
			}
		}

		logger.Warn("origin rejected", slog.String("origin", origin))

		return false
	}
}

func setupRoleManager() *rolemanager.RoleManager {
	methodRoles := map[string][]types.Role{
		"/chat.ConversationService/ConversationNewList": {types.RoleSupport},
		"/chat.ConversationService/ConversationOwnList": {types.RoleSupport},
		"/chat.ConversationService/ChatHistory":         {types.RoleSupport, types.RoleSuperUser, types.RoleClient, types.RoleGuest},
		"/chat.ConversationService/OpenConversation":    {types.RoleSupport},
		"/chat.ConversationService/CloseConversation":   {types.RoleSupport},
	}

	return rolemanager.NewRoleManager(methodRoles)
}
