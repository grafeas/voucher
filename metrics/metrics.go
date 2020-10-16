package metrics

import (
	"time"
)

type Client interface {
	CheckRunStart(string)
	CheckRunLatency(string, time.Duration)
	CheckAttestationLatency(string, time.Duration)
	CheckRunFailure(string)
	CheckRunError(string, error)
	CheckRunSuccess(string)
	CheckAttestationStart(string)
	CheckAttestationError(string, error)
	CheckAttestationSuccess(string)
}
