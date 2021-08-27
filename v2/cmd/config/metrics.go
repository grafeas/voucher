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
	} else if ddApiKey := viper.GetString("statsd.datadog_apikey"); ddApiKey != "" {
		// FIXME: bring back datadog
		// dd := datadog.NewAPIClient(datadog.NewConfiguration())
		// ddAppKey := viper.GetString("stats.datadog.appkey")
		// client = metrics.NewDatadogStatsd(dd, ddApiKey, ddAppKey, tags)
	}
	log.Printf("No metrics client configured")
	return &metrics.NoopClient{}, nil
}
