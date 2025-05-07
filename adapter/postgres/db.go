package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/logger"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type Database struct {
	db         *sql.DB
	mu         sync.Mutex
	statements map[statementKey]*sql.Stmt
}

func New(cfg Config) *Database {
	conn, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		logger.L().Error("Connection to postgres failed", slog.String("error", err.Error()))
	}

	if err = conn.Ping(); err != nil {
		panic(err)
	}

	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	logger.L().Info("postgres connection established successfully!!!")

	return &Database{
		db:         conn,
		statements: make(map[statementKey]*sql.Stmt),
	}
}

func (db *Database) Conn() *sql.DB {
	return db.db
}

func (db *Database) PrepareStatement(ctx context.Context, key statementKey, query string) (*sql.Stmt, error) {
	const op = "repository.postgres.PrepareStatement"

	db.mu.Lock()
	defer db.mu.Unlock()

	if stmt, ok := db.statements[key]; ok {
		return stmt, nil
	}

	stmt, pErr := db.db.PrepareContext(ctx, query) //nolint:sqlclosecheck // this is closed in CloseStatements method.
	if pErr != nil {
		richErr := richerror.New(op).WithWrapError(pErr).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]interface{}{"key": key, "query": query})

		return nil, richErr
	}

	db.statements[key] = stmt

	return stmt, nil
}

func (db *Database) CloseStatements() error {
	const op = "repository.postgres.CloseStatement"

	db.mu.Lock()
	defer db.mu.Unlock()

	for key, stmt := range db.statements {
		cErr := stmt.Close()
		if cErr != nil {
			richErr := richerror.New(op).WithWrapError(cErr).WithMeta(map[string]interface{}{"key": key})

			return richErr
		}
	}

	return nil
}

func (db *Database) Close() error {
	return db.db.Close()
}
