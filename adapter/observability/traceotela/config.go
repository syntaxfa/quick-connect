package traceotela

import (
	"crypto/tls"
	"time"
)

type Config struct {
	Mode                   Mode          `koanf:"mode"`
	HTTPCollectionEndpoint string        `koanf:"http_collection_endpoint"`
	GRPCCollectionEndpoint string        `koanf:"grpc_collection_endpoint"`
	BatchTimeout           time.Duration `koanf:"batch_timeout"`
	BatchSize              int           `koanf:"batch_size"`
	SSLMode                bool          `koanf:"ssl_mode"`
	TLSConfig              *tls.Config   `koanf:"tls_config"`
	SamplingRatio          float64       `koanf:"sampling_ratio"`
}

type Mode string

const (
	ModeDisable = "disable"
	ModeConsole = "console"
	ModeRemote  = "remote"
)
