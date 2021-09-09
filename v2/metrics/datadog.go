package metrics

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/url"
	"sort"
	"strings"
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
	// map of timingKey() to timestamp, to values
	durationData map[string]map[float64][]float64
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
		authCtx:      context.WithValue(context.Background(), datadog.ContextAPIKeys, keys),
		metrics:      client.MetricsApi,
		events:       client.EventsApi,
		now:          func() float64 { return float64(time.Now().Unix()) },
		durationData: make(map[string]map[float64][]float64),
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
	tk := timingKey(metric, tags)

	d.mu.Lock()
	defer d.mu.Unlock()

	if existing := d.findSeries(metric, tags); existing != nil {
		d.durationData[tk][now] = append(d.durationData[tk][now], val)
		return nil
	}

	// Not found, create
	s := datadog.NewSeries(metric, nil)
	s.SetType(durationType)
	s.SetTags(tags)
	s.SetInterval(1)
	d.series = append(d.series, s)
	d.durationData[tk] = map[float64][]float64{now: {val}}
	return nil
}

func timingKey(metric string, tags []string) string {
	return fmt.Sprintf("%s %s", metric, strings.Join(tags, ","))
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
	series := d.flushSeries()
	if len(series) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(d.authCtx, submitTimeout)
	defer cancel()

	// TODO: this is sync, that is not great
	if _, _, err := d.metrics.SubmitMetrics(ctx, *datadog.NewMetricsPayload(series)); err != nil {
		log.Println("error submitting metrics to datadog", err)
	}
}

func (d *datadogStatsd) flushSeries() []datadog.Series {
	d.mu.Lock()
	defer d.mu.Unlock()

	series := make([]datadog.Series, 0, len(d.series))
	for _, s := range d.series {
		if s.GetType() != durationType {
			series = append(series, *s)
			continue
		}

		// Construct series from captured samples:
		tk := timingKey(s.GetMetric(), s.GetTags())
		var counts, averages, p95s, p99s, maxes [][]float64
		for ts, data := range d.durationData[tk] {
			sort.Float64s(data)
			var sum float64
			for _, f := range data {
				sum += f
			}
			count := float64(len(data))
			counts = append(counts, []float64{ts, count})
			averages = append(averages, []float64{ts, sum / count})
			p95s = append(p95s, []float64{ts, percentile(data, 0.95)})
			p99s = append(p99s, []float64{ts, percentile(data, 0.99)})
			maxes = append(maxes, []float64{ts, data[len(data)-1]})
		}
		series = append(series, cloneSeries(s, "count", counts))
		series = append(series, cloneSeries(s, "avg", averages))
		series = append(series, cloneSeries(s, "95percentile", p95s))
		series = append(series, cloneSeries(s, "99percentile", p99s))
		series = append(series, cloneSeries(s, "max", maxes))
	}
	d.series = nil
	d.durationData = make(map[string]map[float64][]float64)
	return series
}

func cloneSeries(s *datadog.Series, suffix string, points [][]float64) datadog.Series {
	return datadog.Series{
		Host:           s.Host,
		Interval:       s.Interval,
		Metric:         fmt.Sprintf("%s.%s", s.GetMetric(), suffix),
		Points:         points,
		Tags:           s.Tags,
		Type:           s.Type,
		UnparsedObject: s.UnparsedObject,
	}
}

func percentile(data []float64, p float64) float64 {
	pos := float64(len(data)) * p
	if math.Round(pos) == pos {
		// Return exact value at percentile
		return data[int(pos)]
	}

	return (data[int(pos-1)] + data[int(pos)]) / 2
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
