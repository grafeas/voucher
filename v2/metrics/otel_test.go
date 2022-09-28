package metrics_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/grafeas/voucher/v2/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestOpenTelemetryClient(t *testing.T) {
	reader := metric.NewManualReader()

	// Exercise the client to produce some metrics:
	client, err := metrics.NewOpenTelemetryClient(metric.NewMeterProvider(metric.WithReader(reader)))
	require.NoError(t, err)
	for i := 0; i < 10; i++ {
		client.CheckRunStart("diy")
		client.CheckAttestationLatency("diy", time.Duration(rand.Intn(500)+500)*time.Millisecond)
	}

	metrics, err := reader.Collect(context.Background())
	require.NoError(t, err)
	require.Len(t, metrics.ScopeMetrics, 1)
	require.Len(t, metrics.ScopeMetrics[0].Metrics, 11, "total metric count")

	// Verify the metrics we triggered are present:
	names := make(map[string]struct{}, len(metrics.ScopeMetrics[0].Metrics))
	for _, m := range metrics.ScopeMetrics[0].Metrics {
		names[m.Name] = struct{}{}
		switch m.Name {
		case "voucher_check_run_start_total":
			agg := m.Data.(metricdata.Sum[int64])
			assert.Len(t, agg.DataPoints, 1)
			assert.Equal(t, int64(10), agg.DataPoints[0].Value)
		case "voucher_check_attestation_latency_milliseconds":
			agg := m.Data.(metricdata.Histogram)
			assert.Len(t, agg.DataPoints, 1)
			assert.Equal(t, uint64(10), agg.DataPoints[0].Count)
			assert.GreaterOrEqual(t, *agg.DataPoints[0].Min, float64(500))
			assert.LessOrEqual(t, *agg.DataPoints[0].Max, float64(1000))
		}
	}
	assert.Contains(t, names, "voucher_check_run_start_total")
	assert.Contains(t, names, "voucher_check_attestation_latency_milliseconds")
}
