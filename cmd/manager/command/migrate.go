package command

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/pkg/migrator"
)

type Migrate struct {
	cfg    postgres.Config
	logger *slog.Logger
	limit  int
}

func (m Migrate) Command(cfg postgres.Config, logger *slog.Logger) *cobra.Command {
	m.cfg = cfg
	m.logger = logger

	cmd := &cobra.Command{
		Use:       "migrate",
		Short:     "run migrations",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"up", "down"},
		Run: func(_ *cobra.Command, args []string) {
			m.run(args)
		},
	}

	cmd.Flags().IntVar(&m.limit, "limit", 0, "limit the number of migrations to apply (default 0 means unlimited)")

	return cmd
}

func (m Migrate) run(args []string) {
	if len(args) != 1 {
		m.logger.Error("invalid arguments given", slog.Any("args", args))
	}

	mgr := migrator.New(migrator.Config{
		Host:     m.cfg.Host,
		Port:     m.cfg.Port,
		Username: m.cfg.Username,
		Password: m.cfg.Password,
		DBName:   m.cfg.DBName,
		SSLMode:  m.cfg.SSLMode,
	}, m.cfg.PathOfMigration)

	defer func() {
		if err := mgr.Close(); err != nil {
			m.logger.Error("migrator close failed", slog.String("error", err.Error()))
		}
	}()

	limit := m.limit
	if limit < 0 {
		m.logger.Error("invalid limit value", slog.Int("limit", limit))

		return
	}

	switch args[0] {
	case "up":
		if n, err := mgr.Up(limit); err != nil {
			m.logger.Error("error migrations up", slog.String("error", err.Error()))
		} else {
			m.logger.Info(fmt.Sprintf("applied %d migrations!", n), slog.Int("migration_count", n))
		}
	case "down":
		if n, err := mgr.Down(limit); err != nil {
			m.logger.Error("error migrations down", slog.String("error", err.Error()))
		} else {
			m.logger.Info(fmt.Sprintf("downgrade %d migrations!", n), slog.Int("migration_count", n))
		}
	default:
		log.Println("please specify a migrations direction with up or down")
	}

	m.logger.Info(fmt.Sprintf("migrations %s successfully run with CLI", args[0]))
}
