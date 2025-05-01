package time

import "time"

type Provider interface {
	Now() time.Time
}

type timeProvider struct{}

func NewTimeProvider() Provider {
	return timeProvider{}
}

func (t timeProvider) Now() time.Time {
	return time.Now()
}
