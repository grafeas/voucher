package metrics

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

type DogStatsdClient struct {
	client       *statsd.Client
	samplingRate float64
}

func NewDogStatsdClient(addr string, samplingRate float64, tags []string) (*DogStatsdClient, error) {
	client, err := statsd.New(addr, statsd.WithTags(tags))
	if err != nil {
		return nil, err
	}

	return &DogStatsdClient{
		client:       client,
		samplingRate: samplingRate,
	}, nil
}

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

func createDataDogErrorEvent(check, title string, err error) *statsd.Event {
	event := statsd.NewEvent(title, err.Error())
	event.AlertType = statsd.Error
	event.Priority = statsd.Low
	event.Tags = []string{"check:" + check}
	return event
}
