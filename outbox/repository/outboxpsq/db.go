package outboxpsq

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *DB {
	return &DB{
		conn: conn,
	}
}
