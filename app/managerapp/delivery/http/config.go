package http

import "time"

type Config struct {
	RegisterGuestMaxHint       int           `koanf:"register_guest_max_hint"`
	RegisterGuestDurationLimit time.Duration `koanf:"register_guest_duration_limit"`
}
