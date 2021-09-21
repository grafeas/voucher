package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var errNoHost = errors.New("cannot create client with empty hostname")

// Client is a client for the Voucher API.
type Client struct {
	hostname   *url.URL
	httpClient *http.Client
	username   string
	password   string
}

// SetBasicAuth adds the username and password to the Client struct
func (c *Client) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
}

// NewClient creates a new Client set to connect to the passed
// hostname.
func NewClient(hostname string) (*Client, error) {
	var err error

	if "" == hostname {
		return nil, errNoHost
	}

	hostnameURL, err := url.Parse(hostname)
	if nil != err {
		return nil, fmt.Errorf("could not parse voucher hostname: %s", err)
	}

	if "" == hostnameURL.Scheme {
		hostnameURL.Scheme = "https"
	}

	client := &Client{
		hostname:   hostnameURL,
		httpClient: &http.Client{},
	}

	return client, nil
}
