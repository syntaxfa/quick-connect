package command

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/app/notificationapp"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
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
	const op = "cmd.notification.server.run"

	re := redis.New(s.cfg.Redis, s.logger)
	pg := postgres.New(s.cfg.Postgres, s.logger)

	app := notificationapp.Setup(s.cfg, s.logger, s.trap, re, pg)

	app.Start()

	if cErr := re.Close(); cErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), s.logger)
	}
	s.logger.Info("redis connection closed")

	pg.Close()
	s.logger.Info("postgres connection closed")
}
