package server

import (
	"fmt"
	"time"
)

// Config is a structure which contains Server configuration.
type Config struct {
	Port        int
	Timeout     int
	RequireAuth bool
	Username    string
	PassHash    string
}

// Address is the address of the Server.
func (config *Config) Address() string {
	return fmt.Sprintf(":%d", config.Port)
}

// TimeoutDuration returns the configured timeout for this Server.
func (config *Config) TimeoutDuration() time.Duration {
	return time.Duration(config.Timeout) * time.Second
}
