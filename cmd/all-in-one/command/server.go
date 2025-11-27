package command

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/app/adminapp"
	"github.com/syntaxfa/quick-connect/app/chatapp"
	"github.com/syntaxfa/quick-connect/app/managerapp"
	"github.com/syntaxfa/quick-connect/app/notificationapp"
)

type Server struct {
	cfg    ServiceConfig
	logger Logger
}

func (s Server) Command(cfg ServiceConfig, logger Logger, trap <-chan os.Signal) *cobra.Command {
	s.cfg = cfg
	s.logger = logger

	run := func(_ *cobra.Command, _ []string) {
		s.run(trap)
	}

	return &cobra.Command{
		Use:   "start",
		Short: "start quick connect in code-level monolith",
		Run:   run,
	}
}

func (s Server) run(trap <-chan os.Signal) {
	managerTrap := make(chan os.Signal, 1)
	chatTrap := make(chan os.Signal, 1)
	notificationTrap := make(chan os.Signal, 1)
	adminTrap := make(chan os.Signal, 1)

	reFactory := newRedisFactory(s.logger.ManagerLog)
	defer reFactory.closeAll()

	managerPsqAdapter := postgres.New(s.cfg.ManagerCfg.Postgres, s.logger.ManagerLog)
	defer func() {
		managerPsqAdapter.Close()

		s.logger.ManagerLog.Info("manager postgres connection closed")
	}()

	chatPsqAdapter := postgres.New(s.cfg.ChatCfg.Postgres, s.logger.ChatLog)
	defer func() {
		chatPsqAdapter.Close()

		s.logger.ChatLog.Info("chat postgres connection closed")
	}()

	notificationPsqAdapter := postgres.New(s.cfg.NotificationCfg.Postgres, s.logger.NotificationLog)
	defer func() {
		notificationPsqAdapter.Close()

		s.logger.NotificationLog.Info("notification postgres connection closed")
	}()

	var wg sync.WaitGroup

	managerApp, _, _ := managerapp.Setup(s.cfg.ManagerCfg, s.logger.ManagerLog, managerTrap, managerPsqAdapter,
		reFactory.newConnection(s.cfg.ManagerCfg.Redis, s.logger.ManagerLog))

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.ManagerLog.Info("Starting Manager App...")
		managerApp.Start()
		s.logger.ManagerLog.Info("Manager App Stopped")
	}()

	chatApp, _ := chatapp.Setup(s.cfg.ChatCfg, s.logger.ChatLog, chatTrap, chatPsqAdapter,
		reFactory.newConnection(s.cfg.ChatCfg.Redis, s.logger.ChatLog))

	wg.Add(1)
	func() {
		defer wg.Done()
		s.logger.ChatLog.Info("Starting Chat App...")
		chatApp.Start()
		s.logger.ChatLog.Info("Chat App Stopped")
	}()

	notificationApp, _ := notificationapp.Setup(s.cfg.NotificationCfg, s.logger.NotificationLog, notificationTrap,
		reFactory.newConnection(s.cfg.NotificationCfg.Redis, s.logger.NotificationLog), notificationPsqAdapter)

	wg.Add(1)
	func() {
		defer wg.Done()
		s.logger.NotificationLog.Info("Starting Notification App...")
		notificationApp.Start()
		s.logger.NotificationLog.Info("Notification App Stopped")
	}()

	adminApp := adminapp.Setup(s.cfg.AdminCfg, s.logger.AdminLog, adminTrap)

	wg.Add(1)
	func() {
		defer wg.Done()
		s.logger.AdminLog.Info("Starting admin App...")
		adminApp.Start()
		s.logger.AdminLog.Info("Admin App Stopped")
	}()

	sig := <-trap
	s.logger.ManagerLog.Info("Received shutdown signal", slog.String("signal", sig.String()))

	managerTrap <- sig
	chatTrap <- sig
	notificationTrap <- sig
	adminTrap <- sig

	s.logger.ManagerLog.Info("Waiting for services to shut down...")
	wg.Wait()
	s.logger.ManagerLog.Info("All services stopped. Exiting.")
}

type redisFactory struct {
	conns        map[string]*redis.Adapter
	mu           sync.Mutex
	globalLogger *slog.Logger
}

func newRedisFactory(globalLogger *slog.Logger) *redisFactory {
	return &redisFactory{
		conns:        make(map[string]*redis.Adapter),
		globalLogger: globalLogger,
	}
}

// newConnection if connection exists, returns.
func (r *redisFactory) newConnection(cfg redis.Config, logger *slog.Logger) *redis.Adapter {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, ok := r.conns[fmt.Sprintf("%d@%s:%d/%s", cfg.DB, cfg.Host, cfg.Port, cfg.Password)]
	if !ok {
		conn = redis.New(cfg, logger)

		logger.Info("create redis connection")
	} else {
		logger.Info("using same redis connection")
	}

	return conn
}

func (r *redisFactory) closeAll() {
	for key, conn := range r.conns {
		if cErr := conn.Close(); cErr != nil {
			r.globalLogger.Error("redis connection closed failed", slog.String("error", cErr.Error()),
				slog.String("connection", key))
		}
	}

	r.globalLogger.Info("redis connections gracefully shutdown")
}
