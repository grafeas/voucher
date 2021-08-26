package config

import (
	"log"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/grafeas/voucher/v2/metrics"
	"github.com/spf13/viper"
)

func MetricsClient() (metrics.Client, error) {
	tags := viper.GetStringSlice("statsd.tags")
	var client metrics.StatsdClient
	if statsdAddr := viper.GetString("statsd.addr"); statsdAddr != "" {
		var err error
		client, err = statsd.New(statsdAddr, statsd.WithTags(tags))
		if err != nil {
			return nil, err
		}
	} else if ddApiKey := viper.GetString("stats.datadog.apikey"); ddApiKey != "" {
		dd := datadog.NewAPIClient(datadog.NewConfiguration())
		ddAppKey := viper.GetString("stats.datadog.appkey")
		client = metrics.NewDatadogStatsd(dd, ddApiKey, ddAppKey, tags)
	} else {
		log.Printf("No metrics client configured")
		return &metrics.NoopClient{}, nil
	}

	sampleRate := viper.GetFloat64("statsd.sample_rate")
	return metrics.NewStatsdMetricsClient(client, sampleRate)
}
