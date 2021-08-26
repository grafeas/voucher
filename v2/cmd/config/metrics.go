package config

import (
	"log"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/grafeas/voucher/v2/metrics"
	"github.com/spf13/viper"
)

func MetricsClient() (metrics.Client, error) {
	statsdAddr := viper.GetString("statsd.addr")
	if statsdAddr == "" {
		log.Printf("No metrics client configured")
		return &metrics.NoopClient{}, nil
	}

	tags := viper.GetStringSlice("statsd.tags")
	client, err := statsd.New(statsdAddr, statsd.WithTags(tags))
	if err != nil {
		return nil, err
	}

	sampleRate := viper.GetFloat64("statsd.sample_rate")
	return metrics.NewStatsdMetricsClient(client, sampleRate)
}
