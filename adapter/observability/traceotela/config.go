package traceotela

import (
	"time"
)

// Config holds the quick-connect traces system configuration.
type Config struct {
	Mode                   Mode          `koanf:"mode"` // Operating mode for the metrics system
	HTTPCollectionEndpoint string        `koanf:"http_collection_endpoint"`
	GRPCCollectionEndpoint string        `koanf:"grpc_collection_endpoint"`
	BatchTimeout           time.Duration `koanf:"batch_timeout"`
	BatchSize              int           `koanf:"batch_size"`
	SSLMode                bool          `koanf:"ssl_mode"`
	SamplingRatio          float64       `koanf:"sampling_ratio"` // Sampling rate (0-1)
}

// Mode specifies the operational mode of the metrics system.
type Mode string

const (
	ModeDisable = "disable" // Disable trace
	ModeConsole = "console" // trace in console
	ModeRemote  = "remote"  // trace in remote backend
)
