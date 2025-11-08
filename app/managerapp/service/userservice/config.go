package userservice

import "time"

type Config struct {
	UserIDCacheExpiration time.Duration `koanf:"user_id_cache_expiration"`
	PasswordDefaultLength int           `koanf:"password_default_length"`
}
