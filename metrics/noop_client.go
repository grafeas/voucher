package metrics

import (
	"time"
)

type NoopClient struct{}

func (*NoopClient) CheckRunStart(string)                          {}
func (*NoopClient) CheckRunLatency(string, time.Duration)         {}
func (*NoopClient) CheckAttestationLatency(string, time.Duration) {}
func (*NoopClient) CheckRunFailure(string)                        {}
func (*NoopClient) CheckRunError(string, error)                   {}
func (*NoopClient) CheckRunSuccess(string)                        {}
func (*NoopClient) CheckAttestationStart(string)                  {}
func (*NoopClient) CheckAttestationError(string, error)           {}
func (*NoopClient) CheckAttestationSuccess(string)                {}
