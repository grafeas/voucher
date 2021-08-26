package metrics

import (
	"context"
	"log"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/DataDog/datadog-go/statsd"
)

type DatadogStatsd struct {
	authCtx context.Context
	metrics *datadog.MetricsApiService
	events  *datadog.EventsApiService
	timeout time.Duration
	tags    []string
	now     func() float64
}

func NewDatadogStatsd(client *datadog.APIClient, apiKey string, appKey string, tags []string) *DatadogStatsd {
	keys := map[string]datadog.APIKey{
		"apiKeyAuth": {Key: apiKey},
		"appKeyAuth": {Key: appKey},
	}
	return &DatadogStatsd{
		authCtx: context.WithValue(context.Background(), datadog.ContextAPIKeys, keys),
		metrics: client.MetricsApi,
		events:  client.EventsApi,
		tags:    tags,
		timeout: 3 * time.Second,
		now:     func() float64 { return float64(time.Now().Unix()) },
	}
}

const (
	durationType = "gauge"
	countType    = "count"
)

func (d *DatadogStatsd) Incr(metric string, tags []string, _ float64) error {
	s := datadog.NewSeries(metric, [][]float64{{d.now(), 1}})
	s.SetType(countType)
	s.SetTags(append(d.tags, tags...))
	d.submit(*s)
	return nil
}

func (d *DatadogStatsd) Timing(metric string, dur time.Duration, tags []string, _ float64) error {
	s := datadog.NewSeries(metric, [][]float64{{d.now(), float64(dur.Milliseconds())}})
	s.SetType(durationType)
	s.SetTags(append(d.tags, tags...))
	d.submit(*s)
	return nil
}

func (d *DatadogStatsd) Event(e *statsd.Event) error {
	ctx, cancel := context.WithTimeout(d.authCtx, d.timeout)
	defer cancel()

	ddEvent := datadog.NewEventCreateRequest(e.Text, e.Title)
	ddEvent.SetAlertType(datadog.EventAlertType(e.AlertType))
	ddEvent.SetAggregationKey(e.AggregationKey)
	ddEvent.SetPriority(datadog.EventPriority(e.Priority))
	ddEvent.SetTags(append(d.tags, e.Tags...))
	if _, _, err := d.events.CreateEvent(ctx, *ddEvent); err != nil {
		log.Println("error submitting event to datadog", err)
	}
	return nil
}

func (d *DatadogStatsd) submit(series ...datadog.Series) {
	ctx, cancel := context.WithTimeout(d.authCtx, d.timeout)
	defer cancel()

	// TODO: this is not batched, that is not great
	// TODO: this is sync, that is not great
	if _, _, err := d.metrics.SubmitMetrics(ctx, *datadog.NewMetricsPayload(series)); err != nil {
		log.Println("error submitting metrics to datadog", err)
	}
}
