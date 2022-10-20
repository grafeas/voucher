package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/grafeas/voucher/v2/metrics"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func MetricsClient(secrets *Secrets) (metrics.Client, error) {
	// Prefer the [metrics] section, for common config - but fallback to the old [statsd] section
	tags := viper.GetStringSlice("metrics.tags")
	if len(tags) == 0 {
		tags = viper.GetStringSlice("statsd.tags")
	}
	backend := viper.GetString("metrics.backend")
	if backend == "" {
		backend = viper.GetString("statsd.backend")
	}

	switch backend {
	case "statsd", "":
		if statsdAddr := viper.GetString("statsd.addr"); statsdAddr != "" {
			sampleRate := viper.GetFloat64("statsd.sample_rate")
			return metrics.NewStatsdClient(statsdAddr, sampleRate, tags)
		}

		log.Printf("No metrics client configured")
		return &metrics.NoopClient{}, nil

	case "datadog":
		if secrets != nil && secrets.Datadog.APIKey != "" && secrets.Datadog.AppKey != "" {
			return metrics.NewDatadogClient(secrets.Datadog.APIKey, secrets.Datadog.AppKey, metrics.WithDatadogTags(tags)), nil
		}
		return &metrics.NoopClient{}, fmt.Errorf("missing secrets for datadog")

	case "otel", "opentelemetry":
		return otelMetricsClient(tags)

	default:
		return &metrics.NoopClient{}, fmt.Errorf("unknown statsd backend: %s", backend)
	}
}

func otelMetricsClient(tags []string) (*metrics.OpenTelemetryClient, error) {
	ctx := context.Background()
	exporter, err := otelExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating otel exporter: %w", err)
	}
	interval := viper.GetDuration("opentelemetry.interval")
	if interval == 0 {
		interval = time.Minute
	}

	attrs := make([]attribute.KeyValue, 0, len(tags)+1)
	for _, tag := range tags {
		s := strings.SplitN(tag, ":", 2)
		attrs = append(attrs, attribute.String(s[0], s[1]))
	}
	attrs = append(attrs, semconv.ServiceNameKey.String("voucher"))

	res, err := resource.New(ctx, resource.WithAttributes(attrs...))
	if err != nil {
		return nil, fmt.Errorf("creating otel resource: %w", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(interval))),
	)
	return metrics.NewOpenTelemetryClient(mp, exporter)
}

func otelExporter(ctx context.Context) (metric.Exporter, error) {
	insecure := viper.GetBool("opentelemetry.insecure")

	addr := viper.GetString("opentelemetry.addr")
	otelURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parsing otel url: %w", err)
	}

	log := logrus.WithFields(logrus.Fields{
		"otel_addr": addr,
		"insecure":  insecure,
		"scheme":    otelURL.Scheme,
		"host":      otelURL.Host,
	})

	switch otelURL.Scheme {
	case "grpc":
		opts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(otelURL.Host),
		}
		if insecure {
			opts = append(opts, otlpmetricgrpc.WithInsecure())
		}
		log.Info("creating otel exporter")
		return otlpmetricgrpc.New(ctx, opts...)
	case "http", "https":
		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpoint(otelURL.Host),
		}
		if insecure {
			opts = append(opts, otlpmetrichttp.WithInsecure())
		}
		log.Info("creating otel exporter")
		return otlpmetrichttp.New(ctx, opts...)
	default:
		return nil, fmt.Errorf("unknown otel scheme: %s", otelURL.Scheme)
	}
}
