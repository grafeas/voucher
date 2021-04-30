package subscriber

import (
	"time"

	"github.com/grafeas/voucher/v2/server"
)

// Config stores the necessary details for a Subscriber
type Config struct {
	Server         *server.Server
	Project        string
	Subscription   string
	RequiredChecks []string
	DryRun         bool
	Timeout        int
}

// TimeoutDuration returns the configured timeout for this Server.
func (c *Config) TimeoutDuration() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}
