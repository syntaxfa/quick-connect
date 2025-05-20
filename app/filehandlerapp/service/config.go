package service

import (
	"github.com/syntaxfa/quick-connect/adapter/storage/aws"
	"github.com/syntaxfa/quick-connect/adapter/storage/file"
)

type Config struct {
	StorageType        string      `koanf:"storage_type"`
	TempDeleteDuration int         `koanf:"temp_delete_duration"`
	Local              file.Config `koanf:"local"`
	S3                 aws.Config  `koanf:"s3"`
}
