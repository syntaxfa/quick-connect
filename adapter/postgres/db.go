package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(cfg Config, logger *slog.Logger) *Database {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	config, pErr := pgxpool.ParseConfig(connStr)
	if pErr != nil {
		logger.Error("unable ro parse postgres config", slog.String("error", pErr.Error()))

		panic(pErr)
	}

	config.ConnConfig.Tracer = otelpgx.NewTracer()

	config.MaxConns = cfg.MaxOpenConns
	config.MinConns = cfg.MaxIdleConns
	config.MaxConnLifetime = cfg.ConnMaxLifetime

	pool, cErr := pgxpool.NewWithConfig(context.Background(), config)
	if cErr != nil {
		logger.Error("unable ro create connection pool", slog.String("error", cErr.Error()))

		panic(pErr)
	}

	if pErr := pool.Ping(context.Background()); pErr != nil {
		logger.Error("connection with postgres is not establish!", slog.String("error", pErr.Error()))

		panic(pErr)
	}

	if oErr := otelpgx.RecordStats(pool); oErr != nil {
		logger.Error("unable to record database stats", slog.String("error", oErr.Error()))

		panic(oErr)
	}

	return &Database{
		pool: pool,
	}
}

func (db *Database) Close() {
	db.pool.Close()
}

func (db *Database) Conn() *pgxpool.Pool {
	return db.pool
}
