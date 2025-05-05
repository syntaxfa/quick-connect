package command

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/chatapp"
)

type Server struct {
	cfg    chatapp.Config
	logger *slog.Logger
}

func (s Server) Command(cfg chatapp.Config, logger *slog.Logger, trap chan os.Signal) *cobra.Command {
	s.cfg = cfg
	s.logger = logger

	run := func(_ *cobra.Command, _ []string) {
		s.run(trap)
	}

	return &cobra.Command{
		Use:   "start",
		Short: "run http server",
		Run:   run,
	}
}

func (s Server) run(trap <-chan os.Signal) {
	app := chatapp.Setup(s.cfg, s.logger, trap)
	app.Start()
}
