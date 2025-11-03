package tokenservice

import (
	"crypto/ed25519"
	"log/slog"
	"os"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
)

type Config struct {
	PrivateKeyString string `koanf:"private_key_string"`
	PublicKeyString  string `koanf:"public_key_string"`
	privateKey       ed25519.PrivateKey
	publicKey        ed25519.PublicKey
	AccessExpiry     time.Duration `koanf:"access_expiry"`
	RefreshExpiry    time.Duration `koanf:"refresh_expiry"`
	Issuer           string        `koanf:"issuer"`
	AccessAudience   string        `koanf:"access_audience"`
	RefreshAudience  string        `koanf:"refresh_audience"`
}

type Service struct {
	cfg       Config
	logger    *slog.Logger
	validator *jwtvalidator.Validator
}

func New(cfg Config, logger *slog.Logger) Service {
	var lErr error

	cfg.privateKey, cfg.publicKey, lErr = LoadKeyPair(KeyPair{
		PrivateKeyHex: cfg.PrivateKeyString,
		PublicKeyHex:  cfg.PublicKeyString,
	})
	if lErr != nil {
		logger.Error("token keys not valid", slog.String("error", lErr.Error()))

		os.Exit(1)
	}

	validator := jwtvalidator.New(cfg.PublicKeyString, logger)

	return Service{
		cfg:       cfg,
		logger:    logger,
		validator: validator,
	}
}
