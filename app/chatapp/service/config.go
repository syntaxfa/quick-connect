package service

import "time"

type Config struct {
	WriteWait      time.Duration `koanf:"write_wait"`
	PongWait       time.Duration `koanf:"pong_wait"`
	PingPeriod     time.Duration
	MaxMessageSize int `koanf:"max_message_size"`
}
