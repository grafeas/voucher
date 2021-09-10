package metrics

import (
	"context"
	"fmt"
	"io"
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
var _ io.Closer = (*DatadogClient)(nil)

func NewDatadogClient(apiKey, appKey string, opts ...DatadogClientOpt) *DatadogClient {
	cfg := datadog.NewConfiguration()
	ds := newDatadogStatsd(cfg, apiKey, appKey)
	c := &DatadogClient{
		cfg: cfg,
		StatsdClient: StatsdClient{
			client: ds,
		},
	}
	for _, o := range opts {
		o(c)
	}

	// start the submission loop after processing options, so tests can shorten the interval
	submitLoopCtx, cancel := context.WithCancel(context.Background())
	ds.cancelSubmitLoop = cancel
	go ds.submitLoop(submitLoopCtx)
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

func WithDatadogSubmitInterval(dur time.Duration) DatadogClientOpt {
	return func(c *DatadogClient) {
		c.client.(*datadogStatsd).submitInterval = dur
	}
}

func (d *DatadogClient) Close() error {
	d.client.(*datadogStatsd).Close()
	return nil
}

// datadogStatsd is an alternative statsd.Client that aggregates in memory, with periodic submission to Datadog API.
type datadogStatsd struct {
	authCtx          context.Context
	cancelSubmitLoop context.CancelFunc
	submitInterval   time.Duration

	metrics *datadog.MetricsApiService
	events  *datadog.EventsApiService
	tags    []string
	now     func() float64

	// series tracks all metrics, but timingData is aggregated separately to prefilter calculations for summary metrics (e.g. p95)
	mu         sync.Mutex
	series     []*datadog.Series
	timingData map[string]map[float64][]float64
}

var _ statsdClient = (*datadogStatsd)(nil)

const (
	submitTimeout         = 5 * time.Second
	defaultSubmitInterval = 10 * time.Second
)

func newDatadogStatsd(cfg *datadog.Configuration, apiKey, appKey string) *datadogStatsd {
	client := datadog.NewAPIClient(cfg)
	keys := map[string]datadog.APIKey{
		"apiKeyAuth": {Key: apiKey},
		"appKeyAuth": {Key: appKey},
	}
	return &datadogStatsd{
		authCtx:        context.WithValue(context.Background(), datadog.ContextAPIKeys, keys),
		metrics:        client.MetricsApi,
		events:         client.EventsApi,
		now:            func() float64 { return float64(time.Now().Unix()) },
		timingData:     make(map[string]map[float64][]float64),
		submitInterval: defaultSubmitInterval,
	}
}

const (
	timingType = "gauge"
	countType  = "count"
)

func (d *datadogStatsd) Incr(metric string, tags []string, _ float64) error {
	tags = append(d.tags, tags...)
	now := d.now()

	d.mu.Lock()
	defer d.mu.Unlock()

	if existing := d.findSeries(metric, tags); existing != nil {
		// Series exists: increment or add the per-timestamp count
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
	s.SetInterval(int64(d.submitInterval.Seconds()))
	d.series = append(d.series, s)
	return nil
}

func (d *datadogStatsd) Timing(metric string, dur time.Duration, tags []string, _ float64) error {
	tags = append(d.tags, tags...)
	sort.Strings(tags)
	now := d.now()
	val := float64(dur.Milliseconds())
	tk := timingKey(metric, tags)

	d.mu.Lock()
	defer d.mu.Unlock()

	if existing := d.findSeries(metric, tags); existing != nil {
		// Series exists, we only need to aggregate timing data:
		d.timingData[tk][now] = append(d.timingData[tk][now], val)
		return nil
	}

	// Not found, create
	s := datadog.NewSeries(metric, nil)
	s.SetType(timingType)
	s.SetTags(tags)
	d.series = append(d.series, s)
	d.timingData[tk] = map[float64][]float64{now: {val}}
	return nil
}

// timingKey serializes metric+tags as a key for timingData
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

func (d *datadogStatsd) Close() {
	d.cancelSubmitLoop()
	// flush any buffered metrics. this may block on mutex until the submitLoop finishes
	d.submit()
}

func (d *datadogStatsd) submitLoop(ctx context.Context) {
	t := time.NewTicker(d.submitInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			d.submit()
		case <-ctx.Done():
			return
		}
	}
}

func (d *datadogStatsd) submit() {
	series := d.flushSeries()
	if len(series) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(d.authCtx, submitTimeout)
	defer cancel()
	if _, _, err := d.metrics.SubmitMetrics(ctx, *datadog.NewMetricsPayload(series)); err != nil {
		log.Println("error submitting metrics to datadog", err)
	}
}

func (d *datadogStatsd) flushSeries() []datadog.Series {
	d.mu.Lock()
	defer d.mu.Unlock()

	// most series map 1:1, but timing expands to 5 series
	series := make([]datadog.Series, 0, len(d.series)+len(d.timingData)*4)
	for _, s := range d.series {
		if s.GetType() == timingType {
			series = append(series, d.deriveTimingSeries(s)...)
			continue
		}
		series = append(series, *s)
	}
	d.series = nil
	d.timingData = make(map[string]map[float64][]float64)
	return series
}

func (d *datadogStatsd) deriveTimingSeries(s *datadog.Series) []datadog.Series {
	tk := timingKey(s.GetMetric(), s.GetTags())
	var counts, averages, p95s, p99s, maxes [][]float64
	for ts, data := range d.timingData[tk] {
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

	return []datadog.Series{
		cloneSeries(s, "count", counts),
		cloneSeries(s, "avg", averages),
		cloneSeries(s, "95percentile", p95s),
		cloneSeries(s, "99percentile", p99s),
		cloneSeries(s, "max", maxes),
	}
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
		return data[int(pos)-1]
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
