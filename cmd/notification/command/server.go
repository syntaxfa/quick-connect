package command

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/notificationapp"
)

type Server struct {
	cfg    notificationapp.Config
	logger *slog.Logger
	trap   <-chan os.Signal
}

func (s Server) Command(cfg notificationapp.Config, logger *slog.Logger, trap <-chan os.Signal) *cobra.Command {
	s.cfg = cfg
	s.logger = logger
	s.trap = trap

	run := func(_ *cobra.Command, _ []string) {
		s.run()
	}

	return &cobra.Command{
		Use:   "start",
		Short: "start notification application",
		Run:   run,
	}
}

func (s Server) run() {
	app := notificationapp.Setup(s.cfg, s.logger, s.trap)

	app.Start()
}
