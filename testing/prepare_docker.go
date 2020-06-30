package vtesting

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/distribution/reference"
)

// PrepareDockerTest creates a new http.Client and httptest.Server for testing with.
// The new client is created using the voucher tests specific Auth.
func PrepareDockerTest(t *testing.T, ref reference.Named) (*http.Client, *httptest.Server) {
	t.Helper()

	server := NewTestDockerServer(t)

	auth := NewAuth(server)

	client, err := auth.ToClient(context.TODO(), ref)
	if nil != err {
		t.Fatalf("failed to create client for Docker API test: %s", err)
	}

	err = UpdateClient(client, server)
	if nil != err {
		t.Fatalf("failed to update client for Docker API test: %s", err)
	}

	return client, server
}
