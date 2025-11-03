package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"

	_ "github.com/lib/pq" // Postgres driver
	migrate "github.com/rubenv/sql-migrate"
)

type Migrate struct {
	migrations *migrate.FileMigrationSource
	db         *sql.DB
}

func New(cfg Config, path string) Migrate {
	hostPort := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, hostPort, cfg.DBName, cfg.SSLMode)

	db, oErr := sql.Open("postgres", connStr)
	if oErr != nil {
		panic(oErr)
	}

	if pErr := db.PingContext(context.Background()); pErr != nil {
		panic(pErr)
	}

	return Migrate{
		migrations: &migrate.FileMigrationSource{Dir: path},
		db:         db,
	}
}

// Up migrate
// Will apply at most `max` migrations. Pass 0 for no limit (or use Exec).
func (m Migrate) Up(maxM int) (int, error) {
	n, err := migrate.ExecMax(m.db, "postgres", m.migrations, migrate.Up, maxM)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Down migrate
// Will apply at most `max` migrations. Pass 0 for no limit (or use Exec).
func (m Migrate) Down(maxM int) (int, error) {
	n, err := migrate.ExecMax(m.db, "postgres", m.migrations, migrate.Down, maxM)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Close DB connection.
func (m Migrate) Close() error {
	return m.db.Close()
}
