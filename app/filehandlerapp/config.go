package filehandlerapp

import (
	"github.com/syntaxfa/quick-connect/adapter/aws"
	"github.com/syntaxfa/quick-connect/adapter/file"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type Config struct {
	Logger   logger.Config   `koanf:"logger"`
	Postgres postgres.Config `koanf:"postgres"`
	Local    file.Config     `koanf:"local"`
	S3       aws.Config      `koanf:"s3"`
}
