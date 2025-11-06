package http

import "time"

type Config struct {
	RegisterGuestMaxHint       int           `koanf:"register_guest_max_hint"`
	RegisterGuestDurationLimit time.Duration `koanf:"register_guest_duration_limit"`
	UpdateGuestMaxHint         int           `koanf:"update_guest_max_hint"`
	UpdateGuestDurationLimit   time.Duration `koanf:"update_guest_duration_limit"`
}
