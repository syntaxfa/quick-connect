package outbox

import "time"

type RetrialPolicy struct {
	MaxSendAttemptsEnabled bool `koanf:"max_send_attempts_enabled"`
	MaxSendAttempts        int  `koanf:"max_send_attempts"`
}

type Config struct {
	ProcessInterval           time.Duration `koanf:"process_interval"`
	LockCheckerInterval       time.Duration `koanf:"lock_checker_interval"`
	MaxLockTimeDuration       time.Duration `koanf:"max_lock_time_duration"`
	CleanupWorkerInterval     time.Duration `koanf:"cleanup_worker_interval"`
	RetrialPolicy             RetrialPolicy `koanf:"retrial_policy"`
	MessagesRetentionDuration time.Duration `koanf:"messages_retention_duration"`
}
