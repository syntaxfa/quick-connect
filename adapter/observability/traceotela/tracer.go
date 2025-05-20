package traceotela

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

var once sync.Once

// SetTracer sets the tracer with the specified name.
// This function only takes effect on the first call.
func SetTracer(name string) {
	once.Do(func() {
		tracer = otel.Tracer(name)
	})
}

// Tracer returns the current tracer.
// If the tracer has not been set previously, nil pointer error.
func Tracer() trace.Tracer {
	return tracer
}

func ResetTracer() {
	tracer = nil

	once = sync.Once{}
}
