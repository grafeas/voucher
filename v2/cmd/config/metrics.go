package config

import (
	"log"

	"github.com/grafeas/voucher/v2/metrics"
	"github.com/spf13/viper"
)

func MetricsClient() (metrics.Client, error) {
	tags := viper.GetStringSlice("statsd.tags")
	if statsdAddr := viper.GetString("statsd.addr"); statsdAddr != "" {
		sampleRate := viper.GetFloat64("statsd.sample_rate")
		return metrics.NewStatsdClient(statsdAddr, sampleRate, tags)
	} else if ddAPIKey := viper.GetString("statsd.datadog_api_key"); ddAPIKey != "" {
		ddAppKey := viper.GetString("statsd.datadog_app_key")
		return metrics.NewDatadogClient(ddAPIKey, ddAppKey, metrics.WithDatadogTags(tags)), nil
	}
	log.Printf("No metrics client configured")
	return &metrics.NoopClient{}, nil
}
