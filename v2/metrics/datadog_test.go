package metrics_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/grafeas/voucher/v2/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mockAPIKey         = "api-key"
	mockAppKey         = "app-key"
	mockClock          = 1096329600
	testSubmitInterval = 100 * time.Millisecond
)

func TestDatadogClient_Counter(t *testing.T) {
	var p datadog.MetricsPayload
	metrics := newMockedDatadogClient(t, &p)

	metrics.CheckRunStart("diy")
	metrics.Close()

	require.Len(t, p.Series, 1)
	series := p.Series[0]
	assert.Equal(t, "voucher.check.run.start", series.Metric)
	assert.Equal(t, []string{"check:diy"}, series.GetTags())
	assert.Equal(t, "count", series.GetType())
	assert.Equal(t, [][]float64{{mockClock, 1}}, series.GetPoints())
}

func TestDatadogClient_Counter_Aggregate(t *testing.T) {
	var p datadog.MetricsPayload
	metrics := newMockedDatadogClient(t, &p)

	metrics.CheckRunStart("diy")
	metrics.CheckRunStart("diy")
	metrics.Close()

	require.Len(t, p.Series, 1)
	series := p.Series[0]
	assert.Equal(t, "voucher.check.run.start", series.Metric)
	assert.Equal(t, []string{"check:diy"}, series.GetTags())
	assert.Equal(t, "count", series.GetType())
	assert.Equal(t, [][]float64{{mockClock, 2}}, series.GetPoints())
}

func TestDatadogClient_Timing(t *testing.T) {
	var p datadog.MetricsPayload
	metrics := newMockedDatadogClient(t, &p)

	metrics.CheckRunLatency("diy", 123*time.Millisecond)
	metrics.Close()

	require.Len(t, p.Series, 5)
	count := p.Series[0]
	assert.Equal(t, "voucher.check.run.latency.count", count.Metric)
	assert.Equal(t, "gauge", count.GetType())
	assert.Equal(t, [][]float64{{mockClock, 1}}, count.GetPoints())

	avg := p.Series[1]
	assert.Equal(t, "voucher.check.run.latency.avg", avg.Metric)
	assert.Equal(t, "gauge", avg.GetType())
	assert.Equal(t, [][]float64{{mockClock, 123}}, avg.GetPoints())

	p95 := p.Series[2]
	assert.Equal(t, "voucher.check.run.latency.95percentile", p95.Metric)
	assert.Equal(t, "gauge", p95.GetType())
	assert.Equal(t, [][]float64{{mockClock, 123}}, p95.GetPoints())

	p99 := p.Series[3]
	assert.Equal(t, "voucher.check.run.latency.99percentile", p99.Metric)
	assert.Equal(t, "gauge", p99.GetType())
	assert.Equal(t, [][]float64{{mockClock, 123}}, p99.GetPoints())

	max := p.Series[4]
	assert.Equal(t, "voucher.check.run.latency.max", max.Metric)
	assert.Equal(t, "gauge", max.GetType())
	assert.Equal(t, [][]float64{{mockClock, 123}}, max.GetPoints())
}

func TestDatadogClient_Timing_Aggregate(t *testing.T) {
	var p datadog.MetricsPayload
	metrics := newMockedDatadogClient(t, &p)

	for i := 0; i < 100; i++ {
		metrics.CheckRunLatency("diy", time.Duration(i+1)*time.Millisecond)
	}
	metrics.Close()

	require.Len(t, p.Series, 5)
	count := p.Series[0]
	assert.Equal(t, "voucher.check.run.latency.count", count.Metric)
	assert.Equal(t, "gauge", count.GetType())
	assert.Equal(t, [][]float64{{mockClock, 100}}, count.GetPoints())

	avg := p.Series[1]
	assert.Equal(t, "voucher.check.run.latency.avg", avg.Metric)
	assert.Equal(t, "gauge", avg.GetType())
	assert.Equal(t, [][]float64{{mockClock, 50.5}}, avg.GetPoints())

	p95 := p.Series[2]
	assert.Equal(t, "voucher.check.run.latency.95percentile", p95.Metric)
	assert.Equal(t, "gauge", p95.GetType())
	assert.Equal(t, [][]float64{{mockClock, 95}}, p95.GetPoints())

	p99 := p.Series[3]
	assert.Equal(t, "voucher.check.run.latency.99percentile", p99.Metric)
	assert.Equal(t, "gauge", p99.GetType())
	assert.Equal(t, [][]float64{{mockClock, 99}}, p99.GetPoints())

	max := p.Series[4]
	assert.Equal(t, "voucher.check.run.latency.max", max.Metric)
	assert.Equal(t, "gauge", max.GetType())
	assert.Equal(t, [][]float64{{mockClock, 100}}, max.GetPoints())
}

func TestDatadogClient_AsyncFlush(t *testing.T) {
	var p datadog.MetricsPayload
	metrics := newMockedDatadogClient(t, &p)

	metrics.CheckRunStart("diy")
	require.Len(t, p.Series, 0)

	time.Sleep(testSubmitInterval * 2)

	require.Len(t, p.Series, 1)
	series := p.Series[0]
	assert.Equal(t, "voucher.check.run.start", series.Metric)
	assert.Equal(t, []string{"check:diy"}, series.GetTags())
	assert.Equal(t, "count", series.GetType())
	assert.Equal(t, [][]float64{{mockClock, 1}}, series.GetPoints())
}

func TestDatadogClient_Event(t *testing.T) {
	var p datadog.EventCreateRequest
	metrics := newMockedDatadogClient(t, &p)

	metrics.CheckAttestationError("diy", errors.New("kaboom"))

	assert.Equal(t, "kaboom", p.Text)
	assert.Equal(t, "Voucher Check Attestation Error", p.Title)
	assert.Equal(t, datadog.EVENTALERTTYPE_ERROR, p.GetAlertType())
	assert.Equal(t, datadog.EVENTPRIORITY_LOW, p.GetPriority())
	assert.Equal(t, []string{"check:diy"}, p.GetTags())
	assert.Equal(t, int64(mockClock), p.GetDateHappened())
}

func newMockedDatadogClient(t testing.TB, data interface{}) *metrics.DatadogClient {
	mockDatadog := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h := r.Header["Dd-Api-Key"]; len(h) != 1 || h[0] != mockAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	t.Cleanup(mockDatadog.Close)

	u, _ := url.Parse(mockDatadog.URL)
	return metrics.NewDatadogClient(mockAPIKey, mockAppKey,
		metrics.WithDatadogURL(*u),
		metrics.WithDatadogFrozenClock(mockClock),
		metrics.WithDatadogSubmitInterval(testSubmitInterval),
	)
}
