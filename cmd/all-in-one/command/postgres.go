package command

import (
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
)

type postgresAdapter struct {
	managerPsqAd      *postgres.Database
	chatPsqAd         *postgres.Database
	notificationPsqAd *postgres.Database
	logger            Logger
}

func newPostgresAdapter(cfg ServiceConfig, logger Logger) *postgresAdapter {
	return &postgresAdapter{
		managerPsqAd:      postgres.New(cfg.ManagerCfg.Postgres, logger.ManagerLog),
		chatPsqAd:         postgres.New(cfg.ChatCfg.Postgres, logger.ChatLog),
		notificationPsqAd: postgres.New(cfg.NotificationCfg.Postgres, logger.NotificationLog),
		logger:            logger,
	}
}

func (d *postgresAdapter) closeAll() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.managerPsqAd.Close()
		d.logger.ManagerLog.Info("manager postgres connection closed")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.chatPsqAd.Close()
		d.logger.ChatLog.Info("chat postgres connection closed")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.notificationPsqAd.Close()
		d.logger.NotificationLog.Info("notification postgres connection closed")
	}()

	wg.Wait()
}
