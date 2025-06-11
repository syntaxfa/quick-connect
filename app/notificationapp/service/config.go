package service

import "time"

type Config struct {
	UserIDCacheExpiration time.Duration `koanf:"user_id_cache_expiration"`
}
