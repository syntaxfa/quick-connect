package traceotela

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
)

func InjectHTTPHeader(ctx context.Context, headers http.Header) {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	propagator.Inject(ctx, propagation.HeaderCarrier(headers))
}

func ExtractHTTPHeader(ctx context.Context, headers http.Header) context.Context {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	return propagator.Extract(ctx, propagation.HeaderCarrier(headers))
}
