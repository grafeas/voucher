package metrics

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

type DogStatsdClient struct {
	client       StatsdClient
	samplingRate float64
}

func NewDogStatsdClient(client StatsdClient, samplingRate float64) (*DogStatsdClient, error) {
	return &DogStatsdClient{
		client:       client,
		samplingRate: samplingRate,
	}, nil
}

type StatsdClient interface {
	Incr(string, []string, float64) error
	Timing(string, time.Duration, []string, float64) error
	Event(*statsd.Event) error
}

var _ StatsdClient = (*statsd.Client)(nil)

func (d *DogStatsdClient) CheckRunStart(check string) {
	_ = d.client.Incr("voucher.check.run.start", []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckRunLatency(check string, dur time.Duration) {
	_ = d.client.Timing("voucher.check.run.latency", dur, []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckAttestationLatency(check string, dur time.Duration) {
	_ = d.client.Timing("voucher.check.attestation.latency", dur, []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckRunError(check string, err error) {
	_ = d.client.Incr("voucher.check.run.error", []string{"check:" + check}, d.samplingRate)
	_ = d.client.Event(createDataDogErrorEvent(check, "Voucher Check Run Error", err))
}

func (d *DogStatsdClient) CheckRunFailure(check string) {
	_ = d.client.Incr("voucher.check.run.failure", []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckRunSuccess(check string) {
	_ = d.client.Incr("voucher.check.run.success", []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckAttestationStart(check string) {
	_ = d.client.Incr("voucher.check.attestation.start", []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckAttestationSuccess(check string) {
	_ = d.client.Incr("voucher.check.attestation.success", []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckAttestationError(check string, err error) {
	_ = d.client.Incr("voucher.check.attestation.error", []string{"check:" + check}, d.samplingRate)
	_ = d.client.Event(createDataDogErrorEvent(check, "Voucher Check Attestation Error", err))
}

// PubSubMessageReceived tracks the number of messages received from pub/sub
func (d *DogStatsdClient) PubSubMessageReceived() {
	_ = d.client.Incr("auto_voucher.message.received", []string{}, d.samplingRate)
}

// PubSubTotalLatency tracks the time it takes to process a pub/sub message
func (d *DogStatsdClient) PubSubTotalLatency(duration time.Duration) {
	_ = d.client.Timing("auto_voucher.latency", duration, []string{}, d.samplingRate)
}

func createDataDogErrorEvent(check, title string, err error) *statsd.Event {
	event := statsd.NewEvent(title, err.Error())
	event.AlertType = statsd.Error
	event.Priority = statsd.Low
	event.Tags = []string{"check:" + check}
	return event
}
