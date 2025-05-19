package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(cfg Config) *Database {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	config, pErr := pgxpool.ParseConfig(connStr)
	if pErr != nil {
		log.Fatalf("unable to parse config: %s\n", pErr.Error())
	}

	config.ConnConfig.Tracer = otelpgx.NewTracer()

	config.MaxConns = cfg.MaxOpenConns
	config.MinConns = cfg.MaxIdleConns
	config.MaxConnLifetime = cfg.ConnMaxLifetime

	pool, cErr := pgxpool.NewWithConfig(context.Background(), config)
	if cErr != nil {
		log.Fatalf("unable to create connection pool: %s\n", cErr.Error())
	}

	if pErr := pool.Ping(context.Background()); pErr != nil {
		log.Fatalf("connection with postgres is not estabslish!, %s\n", pErr.Error())
	}

	if oErr := otelpgx.RecordStats(pool); oErr != nil {
		log.Fatalf("unable to record database stats, %s\n", oErr.Error())
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
