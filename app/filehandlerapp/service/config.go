package service

import "github.com/syntaxfa/quick-connect/adapter/aws"

type LocalStorageConfig struct {
	BaseDir string `koanf:"base_dir"`
}

type Config struct {
	StorageType        string             `koanf:"storage_type"`
	TempDeleteDuration int                `koanf:"temp_delete_duration"`
	Local              LocalStorageConfig `koanf:"local"`
	S3                 aws.Config         `koanf:"s3"`
}
