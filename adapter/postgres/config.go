package postgres

import "time"

type Config struct {
	Host            string        `koanf:"host"`
	Port            int           `koanf:"port"`
	Username        string        `koanf:"username"`
	Password        string        `koanf:"password"`
	DBName          string        `koanf:"db_name"`
	SSLMode         string        `koanf:"ssl_mode"`
	MaxIdleConns    int32         `koanf:"max_idle_conns"`
	MaxOpenConns    int32         `koanf:"max_open_conns"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
	PathOfMigration string        `koanf:"path_of_migration"`
}
