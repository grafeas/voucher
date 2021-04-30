package clair

import "net/http"

// Config is a structure that tracks Clair specific configuration.
type Config struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UseBasicAuth returns true if connections made using this configuration
// should be using BasicAuth.
func (c *Config) UseBasicAuth() bool {
	return ("" != c.Username && "" != c.Password)
}

// UpdateRequest sets the Authorization header for the passed request.
func (c *Config) UpdateRequest(request *http.Request) {
	request.SetBasicAuth(c.Username, c.Password)
}
