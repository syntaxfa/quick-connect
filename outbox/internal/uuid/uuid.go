package uuid

import "github.com/google/uuid"

type Provider interface {
	NewUUID() uuid.UUID
}

type uuidProvider struct{}

func NewUUIDProvider() Provider {
	return uuidProvider{}
}

func (u uuidProvider) NewUUID() uuid.UUID {
	return uuid.New()
}
