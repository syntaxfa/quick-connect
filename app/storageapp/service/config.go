package service

type Config struct {
	Driver      Driver
	Bucket      string
	MaxFileSize int64 `koanf:"max_file_size"`
}
