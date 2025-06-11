package postgres

import "github.com/syntaxfa/quick-connect/adapter/postgres"

type DB struct {
	conn *postgres.Database
}

func New(conn *postgres.Database) *DB {
	return &DB{
		conn: conn,
	}
}
