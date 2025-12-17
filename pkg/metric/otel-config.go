package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bruno303/go-toolkit/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type Config struct {
	ApplicationName    string
	ApplicationVersion string
	Environment        string
	Enabled            bool
	Port               int
	Path               string
	Log                log.Logger
}

func SetupOTelMetrics(ctx context.Context, cfg Config) (shutdown func(context.Context) error, err error) {
	if !cfg.Enabled {
		cfg.Log.Info(ctx, "metrics disabled")
		return func(ctx context.Context) error { return nil }, nil
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ApplicationName),
			semconv.ServiceVersion(cfg.ApplicationVersion),
			semconv.DeploymentEnvironmentName(cfg.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	SetMeter(NewOtelMeter(meterProvider.Meter(cfg.ApplicationName)))
	shutdown = func(ctx context.Context) error {
		return meterProvider.Shutdown(ctx)
	}

	go serveMetrics(ctx, cfg)

	return shutdown, nil
}

func serveMetrics(ctx context.Context, cfg Config) {
	cfg.Log.Info(ctx, fmt.Sprintf("serving metrics at port %d path %s", cfg.Port, cfg.Path))
	http.Handle(cfg.Path, promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
	if err != nil {
		cfg.Log.Error(ctx, "error serving http", err)
		return
	}
}
