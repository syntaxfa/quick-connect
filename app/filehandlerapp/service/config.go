package service

type S3Config struct {
	BucketName string `koanf:"bucket_name"`
	Region     string `koanf:"region"`
	S3Key      string `koanf:"s3_key"`
}

type LocalStorageConfig struct {
	BaseDir string `koanf:"base_dir"`
}

type Config struct {
	StorageType        string             `koanf:"storage_type"`
	TempDeleteDuration int                `koanf:"temp_delete_duration"`
	Local              LocalStorageConfig `koanf:"local"`
	S3                 S3Config           `koanf:"s3"`
}
