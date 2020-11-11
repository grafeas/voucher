package config

import (
	"log"

	"github.com/grafeas/voucher/metrics"
	"github.com/spf13/viper"
)

func MetricsClient() (metrics.Client, error) {
	statsdAddr := viper.GetString("statsd.addr")
	if statsdAddr != "" {
		sampleRate := viper.GetFloat64("statsd.sample_rate")
		tags := viper.GetStringSlice("statsd.tags")
		log.Printf("Sending metrics to StatsD client: %v", statsdAddr)
		return metrics.NewDogStatsdClient(statsdAddr, sampleRate, tags)
	}

	log.Printf("No metrics client configured")
	return &metrics.NoopClient{}, nil
}
