package metrics

import (
	"time"
)

type NoopClient struct{}

func (*NoopClient) CheckRunLatency(string, time.Duration)         {}
func (*NoopClient) CheckAttestationLatency(string, time.Duration) {}
func (*NoopClient) CheckRunFailure(string)                        {}
func (*NoopClient) CheckRunError(string)                          {}
func (*NoopClient) CheckAttestationError(string)                  {}
