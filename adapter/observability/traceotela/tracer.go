package traceotela

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

var once sync.Once

func Tracer(name string) trace.Tracer {
	once.Do(func() {
		tracer = otel.Tracer(name)
	})

	return tracer
}
