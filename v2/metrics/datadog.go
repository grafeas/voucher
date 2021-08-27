package metrics

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/DataDog/datadog-go/statsd"
)

// DatadogClient emits metrics directly to Datadog.
type DatadogClient struct {
	StatsdClient
	cfg *datadog.Configuration
}

func NewDatadogClient(apiKey, appKey string, opts ...DatadogClientOpt) *DatadogClient {
	cfg := datadog.NewConfiguration()
	c := &DatadogClient{
		cfg: cfg,
		StatsdClient: StatsdClient{
			client: newDatadogStatsd(cfg, apiKey, appKey),
		},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type DatadogClientOpt func(*DatadogClient)

func WithDatadogTags(tags []string) DatadogClientOpt {
	return func(c *DatadogClient) {
		c.client.(*datadogStatsd).tags = tags
	}
}

func WithDatadogURL(datadog url.URL) DatadogClientOpt {
	return func(c *DatadogClient) {
		c.cfg.Host = datadog.Host
		c.cfg.Scheme = datadog.Scheme
	}
}

func WithDatadogFrozenClock(frozenTime float64) DatadogClientOpt {
	return func(c *DatadogClient) {
		c.client.(*datadogStatsd).now = func() float64 { return frozenTime }
	}
}

// datadogStatsd is an alternative statsd.Client that transmits directly to Datadog
type datadogStatsd struct {
	authCtx context.Context
	metrics *datadog.MetricsApiService
	events  *datadog.EventsApiService
	tags    []string
	now     func() float64
}

var _ statsdClient = (*datadogStatsd)(nil)

const submitTimeout = 3 * time.Second

func newDatadogStatsd(cfg *datadog.Configuration, apiKey, appKey string) *datadogStatsd {
	client := datadog.NewAPIClient(cfg)
	keys := map[string]datadog.APIKey{
		"apiKeyAuth": {Key: apiKey},
		"appKeyAuth": {Key: appKey},
	}
	return &datadogStatsd{
		authCtx: context.WithValue(context.Background(), datadog.ContextAPIKeys, keys),
		metrics: client.MetricsApi,
		events:  client.EventsApi,
		now:     func() float64 { return float64(time.Now().Unix()) },
	}
}

const (
	durationType = "gauge"
	countType    = "count"
)

func (d *datadogStatsd) Incr(metric string, tags []string, _ float64) error {
	s := datadog.NewSeries(metric, [][]float64{{d.now(), 1}})
	s.SetType(countType)
	s.SetTags(append(d.tags, tags...))
	d.submit(*s)
	return nil
}

func (d *datadogStatsd) Timing(metric string, dur time.Duration, tags []string, _ float64) error {
	s := datadog.NewSeries(metric, [][]float64{{d.now(), float64(dur.Milliseconds())}})
	s.SetType(durationType)
	s.SetTags(append(d.tags, tags...))
	d.submit(*s)
	return nil
}

func (d *datadogStatsd) Event(e *statsd.Event) error {
	ctx, cancel := context.WithTimeout(d.authCtx, submitTimeout)
	defer cancel()

	ddEvent := datadog.NewEventCreateRequest(e.Text, e.Title)
	ddEvent.SetAlertType(datadog.EventAlertType(e.AlertType))
	ddEvent.SetAggregationKey(e.AggregationKey)
	ddEvent.SetPriority(datadog.EventPriority(e.Priority))
	ddEvent.SetTags(append(d.tags, e.Tags...))
	if e.Timestamp.IsZero() {
		ddEvent.SetDateHappened(int64(d.now()))
	} else {
		ddEvent.SetDateHappened(e.Timestamp.Unix())
	}

	if _, _, err := d.events.CreateEvent(ctx, *ddEvent); err != nil {
		log.Println("error submitting event to datadog", err)
	}
	return nil
}

func (d *datadogStatsd) submit(series ...datadog.Series) {
	ctx, cancel := context.WithTimeout(d.authCtx, submitTimeout)
	defer cancel()

	// TODO: this is not batched, that is not great
	// TODO: this is sync, that is not great
	if _, _, err := d.metrics.SubmitMetrics(ctx, *datadog.NewMetricsPayload(series)); err != nil {
		log.Println("error submitting metrics to datadog", err)
	}
}
