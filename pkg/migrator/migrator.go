package migrator

import (
	"context"
	"database/sql"
	"embed"

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
	cfg        Config
}

func NewMigrator(db *sql.DB, migrations embed.FS, cfg Config) Migrator {
	return Migrator{
		db:         db,
		migrations: migrations,
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
	upCommand.Flags().BoolVarP(&opts.ByOne, "by-one", "b", false, "Apply migrations one at a time")

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
		return err
	}

	if opts.ByOne {
		return goose.UpByOneContext(ctx, m.db, m.cfg.Dir)
	}

	if opts.Version != 0 {
		return goose.UpToContext(ctx, m.db, m.cfg.Dir, int64(opts.Version))
	}

	return goose.UpContext(ctx, m.db, m.cfg.Dir)
}

func (m *Migrator) Down(ctx context.Context, opts Options) error {
	goose.SetBaseFS(m.migrations)

	if err := goose.SetDialect(m.cfg.Dialect); err != nil {
		return err
	}

	if opts.Version != 0 {
		return goose.DownToContext(ctx, m.db, m.cfg.Dir, int64(opts.Version))
	}

	return goose.UpContext(ctx, m.db, m.cfg.Dir)
}
