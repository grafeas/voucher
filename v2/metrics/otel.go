package metrics

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
)

// OpenTelemetryClient is a Client using OpenTelemetry metrics.
type OpenTelemetryClient struct {
	checkRunStart   syncint64.Counter
	checkRunFailure syncint64.Counter
	checkRunError   syncint64.Counter
	checkRunSuccess syncint64.Counter
	checkRunLatency syncint64.Histogram

	attestStart   syncint64.Counter
	attestError   syncint64.Counter
	attestSuccess syncint64.Counter
	attestLatency syncint64.Histogram

	pubsubMsgReceived syncint64.Counter
	pubsubMsgLatency  syncint64.Histogram
}

// Please follow https://prometheus.io/docs/practices/naming/ for metric/label naming conventions.

var (
	_             Client = (*OpenTelemetryClient)(nil)
	attrCheckName        = attribute.Key("check_name")
)

// NewOpenTelemetryClient creates a new *OpenTelemetryClient
func NewOpenTelemetryClient(mp metric.MeterProvider) (*OpenTelemetryClient, error) {
	meter := mp.Meter("voucher").SyncInt64()
	client := &OpenTelemetryClient{}
	if err := addRunMetrics(meter, client); err != nil {
		return nil, err
	}
	if err := addAttestMetrics(meter, client); err != nil {
		return nil, err
	}
	if err := addPubSubMetrics(meter, client); err != nil {
		return nil, err
	}
	return client, nil
}

func addRunMetrics(ip syncint64.InstrumentProvider, client *OpenTelemetryClient) (err error) {
	client.checkRunStart, err = ip.Counter("voucher_check_run_start_total")
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_run_start_total counter: %w", err)
	}
	client.checkRunFailure, err = ip.Counter("voucher_check_run_failure_total")
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_run_failure_total counter: %w", err)
	}
	client.checkRunError, err = ip.Counter("voucher_check_run_error_total")
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_run_error_total counter: %w", err)
	}
	client.checkRunSuccess, err = ip.Counter("voucher_check_run_success_total")
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_run_success_total counter: %w", err)
	}
	client.checkRunLatency, err = ip.Histogram("voucher_check_run_latency_milliseconds", instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_run_latency_milliseconds histogram: %w", err)
	}
	return
}

func addAttestMetrics(ip syncint64.InstrumentProvider, client *OpenTelemetryClient) (err error) {
	client.attestStart, err = ip.Counter("voucher_check_attestation_start_total")
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_attestation_start_total counter: %w", err)
	}
	client.attestError, err = ip.Counter("voucher_check_attestation_error_total")
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_attestation_error_total counter: %w", err)
	}
	client.attestSuccess, err = ip.Counter("voucher_check_attestation_success_total")
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_attestation_success_total counter: %w", err)
	}
	client.attestLatency, err = ip.Histogram("voucher_check_attestation_latency_milliseconds", instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		return fmt.Errorf("failed to create voucher_check_attestation_latency_milliseconds histogram: %w", err)
	}
	return
}

func addPubSubMetrics(ip syncint64.InstrumentProvider, client *OpenTelemetryClient) (err error) {
	client.pubsubMsgReceived, err = ip.Counter("voucher_pubsub_message_received_total")
	if err != nil {
		return fmt.Errorf("failed to create voucher_pubsub_message_received_total counter: %w", err)
	}
	client.pubsubMsgLatency, err = ip.Histogram("voucher_pubsub_message_latency_milliseconds", instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		return fmt.Errorf("failed to create voucher_pubsub_message_latency_milliseconds histogram: %w", err)
	}
	return
}

func (o *OpenTelemetryClient) CheckRunStart(check string) {
	o.incr(o.checkRunStart, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) CheckRunFailure(check string) {
	o.incr(o.checkRunFailure, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) CheckRunError(check string, _ error) {
	o.incr(o.checkRunError, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) CheckRunSuccess(check string) {
	o.incr(o.checkRunSuccess, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) CheckAttestationStart(check string) {
	o.incr(o.attestStart, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) CheckAttestationError(check string, _ error) {
	o.incr(o.attestError, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) CheckAttestationSuccess(check string) {
	o.incr(o.attestSuccess, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) CheckRunLatency(check string, dur time.Duration) {
	o.recordSeconds(o.checkRunLatency, dur, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) CheckAttestationLatency(check string, dur time.Duration) {
	o.recordSeconds(o.attestLatency, dur, attrCheckName.String(check))
}

func (o *OpenTelemetryClient) PubSubMessageReceived() {
	o.incr(o.pubsubMsgReceived)
}

func (o *OpenTelemetryClient) PubSubTotalLatency(dur time.Duration) {
	o.recordSeconds(o.pubsubMsgLatency, dur)
}

func (o *OpenTelemetryClient) incr(counter syncint64.Counter, labels ...attribute.KeyValue) {
	o.withContext(func(ctx context.Context) { counter.Add(ctx, 1, labels...) })
}

func (o *OpenTelemetryClient) recordSeconds(hist syncint64.Histogram, dur time.Duration, labels ...attribute.KeyValue) {
	o.withContext(func(ctx context.Context) { hist.Record(ctx, dur.Milliseconds(), labels...) })
}

func (o *OpenTelemetryClient) withContext(f func(context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	f(ctx)
}
