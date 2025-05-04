package migrator

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

type Options struct {
	Version int
	ByOne   bool
}

type Config struct {
	Dialect string
	Dir     string
}

type Migrator struct {
	db         *sql.DB
	migrations embed.FS
	log        *slog.Logger
	cfg        Config
}

func NewMigrator(db *sql.DB, migrations embed.FS, log *slog.Logger, cfg Config) Migrator {
	return Migrator{
		db:         db,
		migrations: migrations,
		log:        log,
		cfg:        cfg,
	}
}

func (m *Migrator) UpCommand() *cobra.Command {
	var opts Options

	upCommand := &cobra.Command{
		Use:   "migrate",
		Short: "Apply migrations to the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Up(cmd.Context(), opts)
		},
	}

	upCommand.Flags().IntVarP(&opts.Version, "version", "v", 0, "Target version to migrate to (0 means all unapplied migrations)")
	upCommand.Flags().BoolVarP(&opts.ByOne, "by-one", "b", false, "Apply the latest migration")

	return upCommand
}

func (m *Migrator) DownCommand() *cobra.Command {
	var opts Options

	upCommand := &cobra.Command{
		Use:   "migrate",
		Short: "Apply migrations to the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Down(cmd.Context(), opts)
		},
	}

	upCommand.Flags().IntVarP(&opts.Version, "version", "v", 0, "Target version to migrate to (0 means all unapplied migrations)")

	return upCommand
}

func (m *Migrator) Up(ctx context.Context, opts Options) error {
	goose.SetBaseFS(m.migrations)

	if err := goose.SetDialect(m.cfg.Dialect); err != nil {
		m.log.Error("invalid dialect", slog.String("error", err.Error()))
		return err
	}

	if opts.ByOne {
		if err := goose.UpByOneContext(ctx, m.db, m.cfg.Dir); err != nil {
			m.log.Error("error at applying migration by one", slog.String("error", err.Error()))
			return err
		}

		m.log.Info("applied migration by one successfully")
		return nil
	}

	if opts.Version != 0 {
		if err := goose.UpToContext(ctx, m.db, m.cfg.Dir, int64(opts.Version)); err != nil {
			m.log.Error(
				"error at applying migration to a specific version",
				slog.String("error", err.Error()),
				slog.Int("version", opts.Version),
			)

			return err
		}

		m.log.Info("applied migration to version successfully", slog.Int("version", opts.Version))
		return nil
	}

	return goose.UpContext(ctx, m.db, m.cfg.Dir)
}

func (m *Migrator) Down(ctx context.Context, opts Options) error {
	goose.SetBaseFS(m.migrations)

	if err := goose.SetDialect(m.cfg.Dialect); err != nil {
		return err
	}

	if opts.Version != 0 {
		if err := goose.DownToContext(ctx, m.db, m.cfg.Dir, int64(opts.Version)); err != nil {
			m.log.Error(
				"error at downgrading migration to a specific version",
				slog.String("error", err.Error()),
				slog.Int("version", opts.Version),
			)

			return err
		}

		m.log.Info("downgraded migration to version successfully", slog.Int("version", opts.Version))
		return nil
	}

	return goose.UpContext(ctx, m.db, m.cfg.Dir)
}
