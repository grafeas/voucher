package github

import "net/http"

// roundTripperWrapper allows us to attach default headers to all githubv4 requests
type roundTripperWrapper struct {
	roundTripper http.RoundTripper
}

// newRoundTripperWrapper creates a new RoundTripperWrapper
func newRoundTripperWrapper(rt http.RoundTripper) *roundTripperWrapper {
	return &roundTripperWrapper{
		roundTripper: rt,
	}
}

// addPreviewSchemaHeaders adds a given array of GitHub API preview schema headers
// to an outgoing request
func addPreviewSchemaHeaders(req *http.Request, previewSchemaHeaders []string) *http.Request {
	for _, header := range previewSchemaHeaders {
		req.Header.Add("Accept", header)
	}

	return req
}

// RoundTrip implements the http RoundTripper interface
func (rtw *roundTripperWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = addPreviewSchemaHeaders(req, previewSchemas)
	return rtw.roundTripper.RoundTrip(req)
}
