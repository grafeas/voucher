package config

import (
	"log"

	"github.com/grafeas/voucher/v2/metrics"
	"github.com/spf13/viper"
)

func MetricsClient(secrets *Secrets) (metrics.Client, error) {
	tags := viper.GetStringSlice("statsd.tags")
	if statsdAddr := viper.GetString("statsd.addr"); statsdAddr != "" {
		sampleRate := viper.GetFloat64("statsd.sample_rate")
		return metrics.NewStatsdClient(statsdAddr, sampleRate, tags)
	}
	if secrets.Datadog.APIKey != "" && secrets.Datadog.AppKey != "" {
		return metrics.NewDatadogClient(secrets.Datadog.APIKey, secrets.Datadog.AppKey, metrics.WithDatadogTags(tags)), nil
	}
	log.Printf("No metrics client configured")
	return &metrics.NoopClient{}, nil
}
