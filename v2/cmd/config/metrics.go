package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/grafeas/voucher/v2/metrics"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func MetricsClient(secrets *Secrets) (metrics.Client, error) {
	tags := viper.GetStringSlice("statsd.tags")

	switch backend := viper.GetString("statsd.backend"); backend {
	case "statsd", "":
		if statsdAddr := viper.GetString("statsd.addr"); statsdAddr != "" {
			sampleRate := viper.GetFloat64("statsd.sample_rate")
			return metrics.NewStatsdClient(statsdAddr, sampleRate, tags)
		}

		log.Printf("No metrics client configured")
		return &metrics.NoopClient{}, nil
	case "otel", "opentelemetry":
		ctx := context.Background()
		exporter, err := otelExporter(ctx)
		if err != nil {
			return nil, fmt.Errorf("creating otel exporter: %w", err)
		}
		interval := viper.GetDuration("statsd.interval")
		if interval == 0 {
			interval = time.Minute
		}
		res, err := resource.New(ctx, resource.WithAttributes(
			semconv.ServiceNameKey.String("voucher"),
		))
		if err != nil {
			return nil, fmt.Errorf("creating otel resource: %w", err)
		}

		mp := metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(interval))),
		)
		return metrics.NewOpenTelemetryClient(mp, exporter)
	case "datadog":
		if secrets.Datadog.APIKey != "" && secrets.Datadog.AppKey != "" {
			return metrics.NewDatadogClient(secrets.Datadog.APIKey, secrets.Datadog.AppKey, metrics.WithDatadogTags(tags)), nil
		}
		return &metrics.NoopClient{}, fmt.Errorf("missing secrets for datadog")
	default:
		return &metrics.NoopClient{}, fmt.Errorf("unknown statsd backend: %s", backend)
	}
}

func otelExporter(ctx context.Context) (metric.Exporter, error) {
	insecure := viper.GetBool("statsd.insecure")

	addr := viper.GetString("statsd.addr")
	otelURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parsing otel url: %w", err)
	}
	switch otelURL.Scheme {
	case "grpc":
		opts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(addr),
		}
		if insecure {
			opts = append(opts, otlpmetricgrpc.WithInsecure())
		}
		return otlpmetricgrpc.New(ctx, opts...)
	case "http", "https":
		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpoint(addr),
		}
		if insecure {
			opts = append(opts, otlpmetrichttp.WithInsecure())
		}
		return otlpmetrichttp.New(ctx, opts...)
	default:
		return nil, fmt.Errorf("unknown otel scheme: %s", otelURL.Scheme)
	}
}
