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
	mockAPIKey = "api-key"
	mockAppKey = "app-key"
)

func TestDatadogClient_Counter(t *testing.T) {
	var p datadog.MetricsPayload
	metrics := newMockedDatadogClient(t, &p)

	metrics.CheckRunStart("diy")

	require.Len(t, p.Series, 1)
	series := p.Series[0]
	assert.Equal(t, "voucher.check.run.start", series.Metric)
	assert.Equal(t, []string{"check:diy"}, series.GetTags())
	assert.Equal(t, "count", series.GetType())
}

func TestDatadogClient_Timing(t *testing.T) {
	var p datadog.MetricsPayload
	metrics := newMockedDatadogClient(t, &p)

	metrics.CheckRunLatency("diy", 123*time.Millisecond)

	require.Len(t, p.Series, 1)
	series := p.Series[0]
	assert.Equal(t, "voucher.check.run.latency", series.Metric)
	assert.Equal(t, "gauge", series.GetType())
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
	return metrics.NewDatadogClient(mockAPIKey, mockAppKey, metrics.WithDatadogURL(*u))
}
