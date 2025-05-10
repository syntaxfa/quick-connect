package command

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/managerapp"
)

type Server struct {
	cfg    managerapp.Config
	logger *slog.Logger
}

func (s Server) Command(cfg managerapp.Config, logger *slog.Logger, trap chan os.Signal) *cobra.Command {
	s.cfg = cfg
	s.logger = logger

	run := func(_ *cobra.Command, _ []string) {
		s.run(trap)
	}

	return &cobra.Command{
		Use:   "start",
		Short: "Start manager application",
		Run:   run,
	}
}

func (s Server) run(trap chan os.Signal) {
	app := managerapp.Setup(s.cfg, s.logger, trap)
	app.Start()
}
