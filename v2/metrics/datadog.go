package metrics

import (
	"context"
	"log"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/DataDog/datadog-go/statsd"
)

// DatadogClient is a metrics.Client that emits directly to Datadog.
type DatadogClient struct {
	StatsdClient
	cfg *datadog.Configuration
}

var _ Client = (*DatadogClient)(nil)

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

func (d *DatadogClient) Close() {
	d.client.(*datadogStatsd).submit()
}

// datadogStatsd is an alternative statsd.Client that transmits directly to Datadog
type datadogStatsd struct {
	authCtx context.Context
	metrics *datadog.MetricsApiService
	events  *datadog.EventsApiService
	tags    []string
	now     func() float64

	mu     sync.Mutex
	series []*datadog.Series
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
	tags = append(d.tags, tags...)
	now := d.now()

	d.mu.Lock()
	defer d.mu.Unlock()

	if existing := d.findSeries(metric, tags); existing != nil {
		for i, p := range existing.Points {
			if p[0] == now {
				existing.Points[i][1]++
				return nil
			}
		}
		existing.Points = append(existing.Points, []float64{now, 1})
		return nil
	}

	// Not found, create
	s := datadog.NewSeries(metric, [][]float64{{now, 1}})
	s.SetType(countType)
	s.SetTags(tags)
	d.series = append(d.series, s)
	return nil
}

func (d *datadogStatsd) Timing(metric string, dur time.Duration, tags []string, _ float64) error {
	tags = append(d.tags, tags...)
	now := d.now()
	val := float64(dur.Milliseconds())

	d.mu.Lock()
	defer d.mu.Unlock()

	if existing := d.findSeries(metric, tags); existing != nil {
		existing.Points = append(existing.Points, []float64{now, val})
		return nil
	}

	// Not found, create
	s := datadog.NewSeries(metric, [][]float64{{now, val}})
	s.SetType(durationType)
	s.SetTags(tags)
	d.series = append(d.series, s)
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

func (d *datadogStatsd) submit() {
	ctx, cancel := context.WithTimeout(d.authCtx, submitTimeout)
	defer cancel()

	d.mu.Lock()
	defer d.mu.Unlock()

	seriesCount := len(d.series)
	if seriesCount == 0 {
		return
	}
	series := make([]datadog.Series, 0, seriesCount)
	for _, s := range d.series {
		series = append(series, *s)
	}
	d.series = nil

	// TODO: this is sync, that is not great
	if _, _, err := d.metrics.SubmitMetrics(ctx, *datadog.NewMetricsPayload(series)); err != nil {
		log.Println("error submitting metrics to datadog", err)
	}
}

func (d *datadogStatsd) findSeries(metric string, tags []string) *datadog.Series {
	sort.Strings(tags)
seriesLoop:
	for _, s := range d.series {
		if s.GetMetric() != metric {
			continue
		}
		sTags := s.GetTags()
		if len(sTags) != len(tags) {
			continue
		}
		sort.Strings(sTags)
		for i, t := range sTags {
			if tags[i] != t {
				continue seriesLoop
			}
		}
		return s
	}
	return nil
}
