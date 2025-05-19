package traceotela

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

func InjectHTTPHeader(ctx context.Context, headers http.Header) {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	propagator.Inject(ctx, propagation.HeaderCarrier(headers))
}

func ExtractHTTPHeader(ctx context.Context, headers http.Header) context.Context {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	return propagator.Extract(ctx, propagation.HeaderCarrier(headers))
}

func InjectGRPCMetadata(ctx context.Context, md *metadata.MD) {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	carrier := metadataCarrier{*md}
	propagator.Inject(ctx, carrier)
}

func ExtractGRPCMetadata(ctx context.Context, md metadata.MD) context.Context {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	carrier := metadataCarrier{md}

	return propagator.Extract(ctx, carrier)
}
