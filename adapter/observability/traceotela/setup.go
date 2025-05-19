package traceotela

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func validateConfig(cfg Config) error {
	if cfg.Mode == ModeRemote && cfg.GRPCCollectionEndpoint == "" && cfg.HTTPCollectionEndpoint == "" {
		return errors.New("http collection endpoint is required for remote mode")
	}

	if cfg.BatchTimeout <= 0 {
		return errors.New("batch timeout must be positive")
	}

	return nil
}

func InitTracer(ctx context.Context, cfg Config, resource *resource.Resource) (func(context.Context) error, error) {
	var provider *trace.TracerProvider

	if vErr := validateConfig(cfg); vErr != nil {
		return nil, vErr
	}

	switch cfg.Mode {
	case ModeConsole:
		traceExporter, sErr := stdouttrace.New(
			stdouttrace.WithPrettyPrint())
		if sErr != nil {
			return nil, sErr
		}

		provider = trace.NewTracerProvider(
			trace.WithBatcher(traceExporter, trace.WithBatchTimeout(cfg.BatchTimeout)),
			trace.WithResource(resource),
		)
	case ModeRemote:
		var traceExporter *otlptrace.Exporter
		var rErr error

		if cfg.GRPCCollectionEndpoint != "" {
			var opts []otlptracegrpc.Option
			opts = append(opts, otlptracegrpc.WithEndpoint(cfg.GRPCCollectionEndpoint))
			if !cfg.SSLMode {
				opts = append(opts, otlptracegrpc.WithInsecure())
			}
			traceExporter, rErr = otlptracegrpc.New(ctx, opts...)
		} else {
			var opts []otlptracehttp.Option
			opts = append(opts, otlptracehttp.WithEndpoint(cfg.HTTPCollectionEndpoint))
			if !cfg.SSLMode {
				opts = append(opts, otlptracehttp.WithInsecure())
			}
			traceExporter, rErr = otlptracehttp.New(ctx, opts...)
		}

		if rErr != nil {
			return nil, rErr
		}

		var sampler trace.Sampler
		switch {
		case cfg.SamplingRatio <= 0:
			sampler = trace.NeverSample()
		case cfg.SamplingRatio <= 1:
			sampler = trace.AlwaysSample()
		default:
			sampler = trace.TraceIDRatioBased(cfg.SamplingRatio)
		}

		provider = trace.NewTracerProvider(
			trace.WithBatcher(traceExporter,
				trace.WithBatchTimeout(cfg.BatchTimeout),
				trace.WithMaxExportBatchSize(cfg.BatchSize),
			),
			trace.WithResource(resource),
			trace.WithSampler(sampler),
		)
	default:
		otel.SetTracerProvider(noop.TracerProvider{})

		return func(_ context.Context) error {
			return nil
		}, nil
	}

	otel.SetTracerProvider(provider)

	return provider.Shutdown, nil
}
