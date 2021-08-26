package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/grafeas/voucher/v2/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mockApiKey = "api-key"
	mockAppKey = "app-key"
)

func TestDatadogStatsd_Incr(t *testing.T) {
	var p datadog.MetricsPayload
	stats := newMockedDatadogStatsd(t, "/api/v1/series", &p)

	stats.Incr("test.counter", []string{"awesome:true"}, 1)

	require.Len(t, p.Series, 1)
	series := p.Series[0]
	assert.Equal(t, "test.counter", series.Metric)
	assert.Equal(t, "count", series.GetType())
}

func TestDatadogStatsd_Timing(t *testing.T) {
	var p datadog.MetricsPayload
	stats := newMockedDatadogStatsd(t, "/api/v1/series", &p)

	stats.Timing("test.duration", 123*time.Millisecond, []string{"awesome:true"}, 1)

	require.Len(t, p.Series, 1)
	series := p.Series[0]
	assert.Equal(t, "test.duration", series.Metric)
	assert.Equal(t, "gauge", series.GetType())
}

func TestDatadogStatsd_Event(t *testing.T) {
	var p datadog.EventCreateRequest
	stats := newMockedDatadogStatsd(t, "/api/v1/events", &p)

	e := statsd.NewEvent("error in check", "kaboom")
	e.AlertType = statsd.Error
	e.Priority = statsd.Low
	e.Tags = []string{"check:test"}
	stats.Event(e)

	assert.Equal(t, "kaboom", p.Text)
	assert.Equal(t, "error in check", p.Title)
}

func newMockedDatadogStatsd(t testing.TB, path string, data interface{}) *metrics.DatadogStatsd {
	mockDatadog := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h := r.Header["Dd-Api-Key"]; len(h) != 1 || h[0] != mockApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.URL.Path != path {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	t.Cleanup(mockDatadog.Close)

	cfg := datadog.NewConfiguration()
	u, _ := url.Parse(mockDatadog.URL)
	cfg.Host = u.Host
	cfg.Scheme = u.Scheme
	dd := datadog.NewAPIClient(cfg)
	return metrics.NewDatadogStatsd(dd, mockApiKey, mockAppKey, nil)
}
