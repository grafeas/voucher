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

func (d *DogStatsdClient) CheckRunLatency(check string, dur time.Duration) {
	d.client.Timing("voucher.check.run.latency", dur, []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckAttestationLatency(check string, dur time.Duration) {
	d.client.Timing("voucher.check.attestation.latency", dur, []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckRunError(check string) {
	d.client.Incr("voucher.check.run.error", []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckRunFailure(check string) {
	d.client.Incr("voucher.check.run.failure", []string{"check:" + check}, d.samplingRate)
}

func (d *DogStatsdClient) CheckAttestationError(check string) {
	d.client.Incr("voucher.check.attestation.error", []string{"check:" + check}, d.samplingRate)
}
