package service

import "time"

type Config struct {
	WriteWait                time.Duration `koanf:"write_wait"`
	PongWait                 time.Duration `koanf:"pong_wait"`
	PingPeriod               time.Duration
	MaxMessageSize           int    `koanf:"max_message_size"`
	ChatChannelName          string `koanf:"chat_channel_name"` // The single channel name for all chat messages
	MessageSnippetCharNumber int    `koanf:"message_snippet_char_number"`
}
