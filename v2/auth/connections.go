package auth

import (
	"errors"
	"net"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// DefaultTransport is a custom implementation of the DefaultTransport.
// It limits the IdleConnTimeout to 10 seconds instead of 90 seconds.
//
// from "net/http/transport.go"
var DefaultTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       10 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// ErrCannotUpdateIdleConnTimeout is an error returned when the Transport
// is unable to be updated.
var ErrCannotUpdateIdleConnTimeout = errors.New("cannot update transport")

// UpdateIdleConnectionsTimeout limits the default timeout for idle connections
// to 10 seconds. If the OAuth2 transport is nil, redefine it with our own
// DefaultTransport. This should limit the number of idle connections left
// open after the initial requests are made, which in turn should reduce the
// chances of the number of available connections being depleted.
func UpdateIdleConnectionsTimeout(client *http.Client) error {
	httpTransport, ok := client.Transport.(*http.Transport)
	if ok {
		httpTransport.IdleConnTimeout = 10 * time.Second
		return nil
	}

	oauth2Transport, ok := client.Transport.(*oauth2.Transport)
	if ok && oauth2Transport.Base == nil {
		oauth2Transport.Base = DefaultTransport

		return nil
	}

	return ErrCannotUpdateIdleConnTimeout
}
