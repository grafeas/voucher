package config_test

import (
	"strings"
	"testing"

	"github.com/grafeas/voucher/v2/cmd/config"
	"github.com/grafeas/voucher/v2/metrics"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsClient(t *testing.T) {
	cases := map[string]struct {
		config       string
		expectedType metrics.Client
	}{
		"disabled": {
			expectedType: &metrics.NoopClient{},
		},
		"statsd from [statsd]": {
			config: `
[statsd]
backend = "statsd"
addr = "localhost:8125"
`,
			expectedType: &metrics.StatsdClient{},
		},
		"statsd from [metrics]": {
			config: `
[metrics]
backend = "statsd"
[statsd]
addr = "localhost:8125"
`,
			expectedType: &metrics.StatsdClient{},
		},
		"datadog from [statsd]": {
			config: `
[statsd]
backend = "datadog"
`,
			expectedType: &metrics.DatadogClient{},
		},
		"datadog from [metrics]": {
			config: `
[metrics]
backend = "datadog"
`,
			expectedType: &metrics.DatadogClient{},
		},
		"otel from [metrics]": {
			config: `
[metrics]
backend = "opentelemetry"
[opentelemetry]
addr = "http://localhost:4317"
`,
			expectedType: &metrics.OpenTelemetryClient{},
		},
	}

	for label, tc := range cases {
		t.Run(label, func(t *testing.T) {
			viper.Reset()
			viper.SetConfigType("toml")
			err := viper.ReadConfig(strings.NewReader(tc.config))
			require.NoError(t, err)

			client, err := config.MetricsClient(&config.Secrets{
				Datadog: config.DatadogSecrets{
					APIKey: "api-key",
					AppKey: "app-key",
				},
			})
			require.NoError(t, err)
			assert.IsType(t, tc.expectedType, client)
		})
	}
}
