package metrics

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// statsdClient is the functionality of github.com/DataDog/datadog-go/statsd used by this implementation
type statsdClient interface {
	Incr(string, []string, float64) error
	Timing(string, time.Duration, []string, float64) error
	Event(*statsd.Event) error
}

var _ statsdClient = (*statsd.Client)(nil)

// StatsdClient is a metrics.Client that emits via statsd UDP protocol.
type StatsdClient struct {
	client       statsdClient
	samplingRate float64
}

var _ Client = (*StatsdClient)(nil)

// NewStatsdClient creates a client to emit metrics to a statsd server, probably dogstatsd.
func NewStatsdClient(addr string, samplingRate float64, tags []string) (*StatsdClient, error) {
	client, err := statsd.New(addr, statsd.WithTags(tags))
	if err != nil {
		return nil, err
	}
	return &StatsdClient{
		client:       client,
		samplingRate: samplingRate,
	}, nil
}

func (d *StatsdClient) CheckRunStart(check string) {
	_ = d.client.Incr("voucher.check.run.start", []string{"check:" + check}, d.samplingRate)
}

func (d *StatsdClient) CheckRunLatency(check string, dur time.Duration) {
	_ = d.client.Timing("voucher.check.run.latency", dur, []string{"check:" + check}, d.samplingRate)
}

func (d *StatsdClient) CheckAttestationLatency(check string, dur time.Duration) {
	_ = d.client.Timing("voucher.check.attestation.latency", dur, []string{"check:" + check}, d.samplingRate)
}

func (d *StatsdClient) CheckRunError(check string, err error) {
	_ = d.client.Incr("voucher.check.run.error", []string{"check:" + check}, d.samplingRate)
	_ = d.client.Event(createDataDogErrorEvent(check, "Voucher Check Run Error", err))
}

func (d *StatsdClient) CheckRunFailure(check string) {
	_ = d.client.Incr("voucher.check.run.failure", []string{"check:" + check}, d.samplingRate)
}

func (d *StatsdClient) CheckRunSuccess(check string) {
	_ = d.client.Incr("voucher.check.run.success", []string{"check:" + check}, d.samplingRate)
}

func (d *StatsdClient) CheckAttestationStart(check string) {
	_ = d.client.Incr("voucher.check.attestation.start", []string{"check:" + check}, d.samplingRate)
}

func (d *StatsdClient) CheckAttestationSuccess(check string) {
	_ = d.client.Incr("voucher.check.attestation.success", []string{"check:" + check}, d.samplingRate)
}

func (d *StatsdClient) CheckAttestationError(check string, err error) {
	_ = d.client.Incr("voucher.check.attestation.error", []string{"check:" + check}, d.samplingRate)
	_ = d.client.Event(createDataDogErrorEvent(check, "Voucher Check Attestation Error", err))
}

// PubSubMessageReceived tracks the number of messages received from pub/sub
func (d *StatsdClient) PubSubMessageReceived() {
	_ = d.client.Incr("auto_voucher.message.received", []string{}, d.samplingRate)
}

// PubSubTotalLatency tracks the time it takes to process a pub/sub message
func (d *StatsdClient) PubSubTotalLatency(duration time.Duration) {
	_ = d.client.Timing("auto_voucher.latency", duration, []string{}, d.samplingRate)
}

func createDataDogErrorEvent(check, title string, err error) *statsd.Event {
	event := statsd.NewEvent(title, err.Error())
	event.AlertType = statsd.Error
	event.Priority = statsd.Low
	event.Tags = []string{"check:" + check}
	return event
}
