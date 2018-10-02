package vtesting

import (
	"net/http"
)

// Transport wraps the http.Transport, but overwrites the URL of all requests
// to point to the same path on the passed hostname. This is to enable us to
// test connections to httptest.Server without changing the client code to
// allow us to add ports to registry URLs (illegal in the reference.Reference
// types in the docker registry library).
type Transport struct {
	hostname  string
	transport *http.Transport
}

// RoundTrip implements http.RoundTripper. It executes a HTTP request with
// the Transport's internal http.Transport, but overrides the hostname before
// doing so.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Host = t.hostname
	return t.transport.RoundTrip(req)
}

// NewTransport creates a new Transport which wraps the http.Transport.
// The purpose of this is to allow us to rewrite URLs to connect to our
// test server.
func NewTransport(hostname string, transport *http.Transport) *Transport {
	newTransport := new(Transport)
	newTransport.transport = transport
	newTransport.hostname = hostname
	return newTransport
}
