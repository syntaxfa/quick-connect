package repository

import (
	"context"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
)

type Database struct {
	db *postgres.Database
}

func New(db *postgres.Database) Database {
	return Database{db: db}
}

func (d Database) GetUserByID(ctx context.Context, _ int) {
	d.db.Conn().Ping(ctx)
}
