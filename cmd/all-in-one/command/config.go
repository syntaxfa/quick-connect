package command

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/app/adminapp"
	"github.com/syntaxfa/quick-connect/app/chatapp"
	"github.com/syntaxfa/quick-connect/app/managerapp"
	"github.com/syntaxfa/quick-connect/app/notificationapp"
)

type ServiceConfig struct {
	ManagerCfg      managerapp.Config
	ChatCfg         chatapp.Config
	NotificationCfg notificationapp.Config
	AdminCfg        adminapp.Config
}

type Logger struct {
	ManagerLog      *slog.Logger
	ChatLog         *slog.Logger
	NotificationLog *slog.Logger
	AdminLog        *slog.Logger
}

type PsqAdapter struct {
	ManagerPsqAdapter      *postgres.Database
	ChatPsqAdapter         *postgres.Database
	NotificationPsqAdapter *postgres.Database
}
