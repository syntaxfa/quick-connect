package command

import (
	"log/slog"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/app/adminapp"
	"github.com/syntaxfa/quick-connect/app/chatapp"
	"github.com/syntaxfa/quick-connect/app/managerapp"
	"github.com/syntaxfa/quick-connect/app/notificationapp"
	"github.com/syntaxfa/quick-connect/pkg/translation"
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
	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		panic(tErr)
	}

	trapSvc := setupTrapService()

	reFactory := newRedisFactory(s.logger.ManagerLog)
	defer reFactory.closeAll()

	postgresAd := newPostgresAdapter(s.cfg, s.logger)
	defer postgresAd.closeAll()

	var wg sync.WaitGroup

	managerApp, tokenSvc, userSvc := managerapp.Setup(s.cfg.ManagerCfg, s.logger.ManagerLog, trapSvc.managerTrap, postgresAd.managerPsqAd,
		reFactory.newConnection(s.cfg.ManagerCfg.Redis))

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.ManagerLog.Info("Starting Manager App...")
		managerApp.Start()
		s.logger.ManagerLog.Info("Manager App Stopped")
	}()

	userInternalLocalAd := manager.NewUserInternalLocalAdapter(&userSvc)
	userLocalAd := manager.NewUserLocalAdapter(&userSvc, t, s.logger.ManagerLog)
	authLocalAdapter := manager.NewAuthLocalAdapter(&userSvc, &tokenSvc, t, s.logger.ManagerLog)
	chatApp, _ := chatapp.Setup(s.cfg.ChatCfg, s.logger.ChatLog, trapSvc.chatTrap, postgresAd.chatPsqAd,
		reFactory.newConnection(s.cfg.ChatCfg.Redis), userInternalLocalAd, authLocalAdapter)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.ChatLog.Info("Starting Chat App...")
		chatApp.Start()
		s.logger.ChatLog.Info("Chat App Stopped")
	}()

	notificationApp, _ := notificationapp.Setup(s.cfg.NotificationCfg, s.logger.NotificationLog, trapSvc.notificationTrap,
		reFactory.newConnection(s.cfg.NotificationCfg.Redis), postgresAd.notificationPsqAd)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.NotificationLog.Info("Starting Notification App...")
		notificationApp.Start()
		s.logger.NotificationLog.Info("Notification App Stopped")
	}()

	adminApp := adminapp.Setup(s.cfg.AdminCfg, s.logger.AdminLog, trapSvc.adminTrap, t, authLocalAdapter,
		userLocalAd, nil)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.AdminLog.Info("Starting admin App...")
		adminApp.Start()
		s.logger.AdminLog.Info("Admin App Stopped")
	}()

	sig := <-trap
	s.logger.ManagerLog.Info("Received shutdown signal", slog.String("signal", sig.String()))

	trapSvc.sendSignal(sig)

	s.logger.ManagerLog.Info("Waiting for services to shut down...")
	wg.Wait()
	s.logger.ManagerLog.Info("All services stopped. Exiting.")
}
