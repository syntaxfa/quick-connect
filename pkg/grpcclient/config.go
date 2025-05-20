package grpcclient

type Config struct {
	Host    string `koanf:"host"`
	Port    int    `koanf:"port"`
	SSLMode bool   `koanf:"ssl_mode"`
	UseOtel bool   `koanf:"use_otel"`
}
