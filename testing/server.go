package vtesting

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// UpdateClient updates the passed http.Client's Transport to support the
// passed httptest.Server's server certificate, and to set a new Transport
// which will override the request's hostname with the testing server's
// hostname.
func UpdateClient(client *http.Client, server *httptest.Server) error {
	serverURL, err := url.Parse(server.URL)
	if nil != err {
		return fmt.Errorf("failed to parse server URL: %s", err)
	}

	certificate := server.Certificate()
	if nil == certificate {
		return errors.New("no TLS certificate active")
	}

	certpool := x509.NewCertPool()
	certpool.AddCert(certificate)

	client.Transport = NewTransport(serverURL.Host, &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certpool,
		},
	})

	return nil
}
