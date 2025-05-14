package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres driver
	migrate "github.com/rubenv/sql-migrate"
)

type Migrate struct {
	migrations *migrate.FileMigrationSource
	db         *sql.DB
}

func New(cfg Config, path string) Migrate {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	db, oErr := sql.Open("postgres", connStr)
	if oErr != nil {
		panic(oErr)
	}

	if pErr := db.Ping(); pErr != nil {
		panic(pErr)
	}

	return Migrate{
		migrations: &migrate.FileMigrationSource{Dir: path},
		db:         db,
	}
}

// Up migrate
// Will apply at most `max` migrations. Pass 0 for no limit (or use Exec).
func (m Migrate) Up(maxM int) (n int, err error) {
	n, err = migrate.ExecMax(m.db, "postgres", m.migrations, migrate.Up, maxM)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Down migrate
// Will apply at most `max` migrations. Pass 0 for no limit (or use Exec).
func (m Migrate) Down(maxM int) (n int, err error) {
	n, err = migrate.ExecMax(m.db, "postgres", m.migrations, migrate.Down, maxM)
	if err != nil {
		return 0, err
	}

	return n, nil
}
