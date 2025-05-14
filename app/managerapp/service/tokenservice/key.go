package tokenservice

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func (s Service) GenerateKeyPair() (*KeyPair, error) {
	op := "auth.service.GenerateKeyPair"

	publicKey, privateKey, gErr := ed25519.GenerateKey(rand.Reader)
	if gErr != nil {
		richErr := richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapError(gErr)
		errlog.ErrLog(richErr, s.logger)

		return nil, richErr
	}

	return &KeyPair{
		PrivateKeyHex: hex.EncodeToString(privateKey),
		PublicKeyHex:  hex.EncodeToString(publicKey),
	}, nil
}

func LoadKeyPair(keyPair KeyPair, logger *slog.Logger) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	op := "auth.service.LoadKeyPair"

	privateKeyBytes, dErr := hex.DecodeString(keyPair.PrivateKeyHex)
	if dErr != nil {
		richErr := richerror.New(op).WithWrapError(dErr).WithKind(richerror.KindUnexpected)
		errlog.ErrLog(richErr, logger)

		return nil, nil, richErr
	}

	publicKeyBytes, dErr := hex.DecodeString(keyPair.PublicKeyHex)
	if dErr != nil {
		richErr := richerror.New(op).WithWrapError(dErr).WithKind(richerror.KindUnexpected)
		errlog.ErrLog(richErr, logger)

		return nil, nil, richErr
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		richErr := richerror.New(op).WithMessage(fmt.Sprintf("invalid private key size: expected %d, got %d",
			ed25519.PrivateKeySize, len(privateKeyBytes)))

		return nil, nil, richErr
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		richErr := richerror.New(op).WithMessage(fmt.Sprintf("invalid public key size: expected %d, got %d",
			ed25519.PublicKeySize, len(privateKeyBytes)))

		return nil, nil, richErr
	}

	return ed25519.PrivateKey(privateKeyBytes), ed25519.PublicKey(publicKeyBytes), nil
}
