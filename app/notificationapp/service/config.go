package service

import "time"

type Config struct {
	UserIDCacheExpiration time.Duration `koanf:"user_id_cache_expiration"`
	ChannelName           string        `json:"channel_name"`
	UserConnectionLimit   int           `koanf:"user_connection_limit"`
	WriteWait             time.Duration `koanf:"write_wait"`
	PongWait              time.Duration `koanf:"pong_wait"`
	MaxMessageSize        int           `koanf:"max_message_size"`
	PublishTimeout        time.Duration `koanf:"publish_timeout"`
	PingPeriod            time.Duration
	DefaultUserLanguage   string `koanf:"default_user_language"`
}
