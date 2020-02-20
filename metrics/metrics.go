package metrics

import (
	"time"
)

type Client interface {
	CheckRunLatency(string, time.Duration)
	CheckAttestationLatency(string, time.Duration)
	CheckRunFailure(string)
	CheckRunError(string)
	CheckAttestationError(string)
}
