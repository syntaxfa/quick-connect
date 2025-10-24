package tokenservice

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func (s Service) GenerateKeyPair() (*KeyPair, error) {
	op := "auth.service.GenerateKeyPair"

	publicKey, privateKey, gErr := ed25519.GenerateKey(rand.Reader)
	if gErr != nil {
		return nil, richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapError(gErr)
	}

	return &KeyPair{
		PrivateKeyHex: hex.EncodeToString(privateKey),
		PublicKeyHex:  hex.EncodeToString(publicKey),
	}, nil
}

func LoadKeyPair(keyPair KeyPair) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	op := "auth.service.LoadKeyPair"

	privateKeyBytes, dErr := hex.DecodeString(keyPair.PrivateKeyHex)
	if dErr != nil {
		return nil, nil, richerror.New(op).WithWrapError(dErr).WithKind(richerror.KindUnexpected)
	}

	publicKeyBytes, dErr := hex.DecodeString(keyPair.PublicKeyHex)
	if dErr != nil {
		return nil, nil, richerror.New(op).WithWrapError(dErr).WithKind(richerror.KindUnexpected)
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, nil, richerror.New(op).WithMessage(fmt.Sprintf("invalid private key size: expected %d, got %d",
			ed25519.PrivateKeySize, len(privateKeyBytes)))
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return nil, nil, richerror.New(op).WithMessage(fmt.Sprintf("invalid public key size: expected %d, got %d",
			ed25519.PublicKeySize, len(privateKeyBytes)))
	}

	return ed25519.PrivateKey(privateKeyBytes), ed25519.PublicKey(publicKeyBytes), nil
}

func (s Service) GetPublicKey() string {
	return s.cfg.PublicKeyString
}
