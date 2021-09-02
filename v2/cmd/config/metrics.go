package config

import (
	"fmt"
	"log"

	"github.com/grafeas/voucher/v2/metrics"
	"github.com/spf13/viper"
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
	case "datadog":
		if secrets.Datadog.APIKey != "" && secrets.Datadog.AppKey != "" {
			return metrics.NewDatadogClient(secrets.Datadog.APIKey, secrets.Datadog.AppKey, metrics.WithDatadogTags(tags)), nil
		}
		return &metrics.NoopClient{}, fmt.Errorf("missing secrets for datadog")
	default:
		return &metrics.NoopClient{}, fmt.Errorf("unknown statsd backend: %s", backend)
	}
}
