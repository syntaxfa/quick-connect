package command

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/adminapp"
	"github.com/syntaxfa/quick-connect/app/chatapp"
	"github.com/syntaxfa/quick-connect/app/managerapp"
	"github.com/syntaxfa/quick-connect/app/notificationapp"
)

type Config struct {
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

type Server struct {
	cfg    Config
	logger Logger
}

func (s Server) Command(cfg Config, logger Logger, trap <-chan os.Signal) *cobra.Command {
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

func (s Server) run(_ <-chan os.Signal) {
}
