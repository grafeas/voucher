package metrics

import (
	"time"
)

type NoopClient struct{}

func (_ *NoopClient) CheckRunLatency(string, time.Duration)         {}
func (_ *NoopClient) CheckAttestationLatency(string, time.Duration) {}
func (_ *NoopClient) CheckRunFailure(string)                        {}
func (_ *NoopClient) CheckRunError(string)                          {}
func (_ *NoopClient) CheckAttestationError(string)                  {}
