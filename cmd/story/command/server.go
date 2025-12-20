package command

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/app/storyapp"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Server struct {
	cfg    storyapp.Config
	logger *slog.Logger
}

func (s Server) Command(cfg storyapp.Config, logger *slog.Logger, trap chan os.Signal) *cobra.Command {
	s.cfg = cfg
	s.logger = logger

	run := func(_ *cobra.Command, _ []string) {
		s.run(trap)
	}

	return &cobra.Command{
		Use:   "start",
		Short: "start story application",
		Run:   run,
	}
}

func (s Server) run(trap <-chan os.Signal) {
	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		s.logger.Error("can't initial translation")

		return
	}

	psqAdapter := postgres.New(s.cfg.Postgres, s.logger)
	defer psqAdapter.Close()

	app, _ := storyapp.Setup(s.cfg, s.logger, trap, t, psqAdapter, nil, nil)

	app.Start()
}
