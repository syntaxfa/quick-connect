package aws

type Config struct {
	BucketName string `koanf:"bucket_name"`
	Region     string `koanf:"region"`
	AccessKey  string `koanf:"access_key"`
	SecretKey  string `koanf:"access_key"`
}
