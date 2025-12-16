package aws

import "time"

type Config struct {
	Endpoint        string `koanf:"endpoint"`
	AccessKeyID     string `koanf:"access_key_id"`
	SecretAccessKey string `koanf:"secret_access_key"`
	BucketName      string `koanf:"bucket_name"`
	Region          string `koanf:"region"`
	UseSSL          bool   `koanf:"use_ssl"`
	// UsePathStyle forces the SDK to use path-style addressing (e.g., https://s3.example.com/bucket/key)
	// instead of virtual-hosted-style addressing (e.g., https://bucket.s3.example.com/key).
	// This is required for most S3-compatible providers like MinIO or local setups.
	UsePathStyle bool `koanf:"use_path_style"`
	// SupportObjectACL indicates whether the storage provider supports Object-Level Access Control Lists (ACLs).
	// If set to true, the application will attempt to set 'public-read' ACL for public files.
	// Set this to false for providers (like Liara, ArvanCloud, or specific MinIO setups) that rely solely on
	// Bucket Policies and ignore or forbid the 'x-amz-acl' header.
	SupportObjectACL bool `koanf:"support_object_acl"`
	// PresignPublicExpire defines how long the generated public URL should be valid
	// when SupportObjectACL is set to false.
	// Example values: "168h" (7 days), "24h" (1 day).
	PresignPublicExpire  time.Duration `koanf:"presign_public_expire"`
	PresignPrivateExpire time.Duration `koanf:"presign_private_expire"`
}
