package command

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/adminapp"
)

type Server struct {
	cfg    adminapp.Config
	logger *slog.Logger
}

func (s Server) Command(cfg adminapp.Config, logger *slog.Logger, trap <-chan os.Signal) *cobra.Command {
	s.cfg = cfg
	s.logger = logger

	run := func(_ *cobra.Command, _ []string) {
		s.run(trap)
	}

	return &cobra.Command{
		Use:   "start",
		Short: "start admin application",
		Run:   run,
	}
}

func (s Server) run(trap <-chan os.Signal) {
	app := adminapp.Setup(s.cfg, s.logger, trap)

	app.Start()
}
