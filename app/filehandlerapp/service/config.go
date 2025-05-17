package service

import "github.com/syntaxfa/quick-connect/adapter/file"

type S3Config struct {
	BucketName string `koanf:"bucket_name"`
	Region     string `koanf:"region"`
	S3Key      string `koanf:"s3_key"`
}

type Config struct {
	StorageType        string      `koanf:"storage_type"`
	TempDeleteDuration int         `koanf:"temp_delete_duration"`
	Local              file.Config `koanf:"local"`
	S3                 S3Config    `koanf:"s3"`
}
