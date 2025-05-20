package metricotela

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

var ErrInvalidPort = errors.New("invalid port number, must between 1 and 65535")

// validateConfig validates configuration and sets default values.
func validateConfig(cfg Config) (Config, error) {
	if cfg.Mode == ModePullBase {
		if cfg.PullConfig.Port == 0 {
			cfg.PullConfig.Port = 12330
		} else if cfg.PullConfig.Port < 1 || cfg.PullConfig.Port > 65535 {
			return cfg, ErrInvalidPort
		}

		if cfg.PullConfig.Path == "" {
			cfg.PullConfig.Path = "/metrics"
		}
	}

	return cfg, nil
}

type server struct {
	s   *http.Server
	mux *http.ServeMux
}

func newServer(port int) *server {
	mux := http.NewServeMux()

	return &server{
		s: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           mux,
			ReadHeaderTimeout: time.Second * 60,
		},
		mux: mux,
	}
}

func initPullBaseMetric(cfg Config, resource *resource.Resource, trap <-chan os.Signal, logger *slog.Logger) error {
	exporter, eErr := otelprom.New()
	if eErr != nil {
		return fmt.Errorf("faliled to create prometheus exporte: %w", eErr)
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(resource),
	)

	otel.SetMeterProvider(provider)

	// Enable runtime metrics collection (goroutines, memory, etc.)
	if rErr := runtime.Start(runtime.WithMeterProvider(provider)); rErr != nil {
		return fmt.Errorf("failed to start runtime metrics collection: %w", rErr)
	}

	var promHandler http.Handler
	defaultRegistry, ok := prometheus.DefaultRegisterer.(*prometheus.Registry)
	if ok {
		promHandler = promhttp.HandlerFor(
			defaultRegistry,
			promhttp.HandlerOpts{},
		)
	} else {
		promHandler = promhttp.Handler()
	}

	httpserver := newServer(cfg.PullConfig.Port)
	httpserver.mux.Handle(cfg.PullConfig.Path, promHandler)

	logger.Info("metrics server started",
		slog.String("endpoint", fmt.Sprintf("http://localhost:%d%s", cfg.PullConfig.Port, cfg.PullConfig.Path)),
		slog.String("mode", "pull"),
	)

	serverErrCh := make(chan error, 1)
	go func() {
		if err := httpserver.s.ListenAndServe(); err != nil {
			serverErrCh <- err
		}
	}()

	select {
	case <-trap:
		logger.Info("received shutdown signal, stopping metrics server...")
	case err := <-serverErrCh:
		logger.Error("metrics server error", slog.String("error", err.Error()))
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	if sErr := httpserver.s.Shutdown(ctx); sErr != nil {
		logger.Error("metric http server shutdown error", slog.String("error", sErr.Error()))

		return fmt.Errorf("metrics server shutdown error: %w", sErr)
	}

	logger.Info("metrics server stopped successfully")

	return nil
}

// InitMetric initializes the metrics system based on the provided configuration.
func InitMetric(cfg Config, resource *resource.Resource, trap <-chan os.Signal, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	var vErr error
	cfg, vErr = validateConfig(cfg)
	if vErr != nil {
		logger.Error("invalid metrics configuration", slog.String("error", vErr.Error()))

		return vErr
	}

	switch cfg.Mode {
	case ModePullBase:
		go func() {
			if err := initPullBaseMetric(cfg, resource, trap, logger); err != nil {
				logger.Error(err.Error())
			}
		}()

		return nil
	default:
		logger.Debug("Metrics disabled, using noop provider")
		otel.SetMeterProvider(noop.MeterProvider{})

		return nil
	}
}
