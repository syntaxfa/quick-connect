package metricotela

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter metric.Meter
	once  sync.Once
)

// SetMeter sets the meter with the specified name.
// This function only takes effect on the first call.
func SetMeter(name string) {
	once.Do(func() {
		meter = otel.Meter(name)
	})
}

// Meter returns the current meter.
// If the meter has not been set previously, returns a default meter.
func Meter() metric.Meter {
	return meter
}

func ResetMeter() {
	meter = nil

	once = sync.Once{}
}
