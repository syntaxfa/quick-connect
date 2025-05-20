package metricotela

// Config holds the quick-connect metrics system configuration.
type Config struct {
	Mode       Mode       `koanf:"mode"`        // Operating mode for the metrics system
	PullConfig PullConfig `koanf:"pull_config"` // Configuration for pull mode
}

// Mode specifies the operational mode of the metrics system.
type Mode string

const (
	ModeDisable  Mode = "disable" // Disable metrics
	ModePullBase Mode = "pull"    // Pull-based mode (Prometheus scrape endpoint)
)

// PullConfig holds the pull mode configuration.
type PullConfig struct {
	Path string `koanf:"path"` // Metrics endpoint path (e.g., "/metrics")
	Port int    `koanf:"port"` // HTTP server port (1-65535)
}
