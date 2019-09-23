package docker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/require"

	vtesting "github.com/Shopify/voucher/testing"
)

// PrepareDockerTest creates a new http.Client and httptest.Server for testing with.
// The new client is created using the voucher tests specific Auth.
func PrepareDockerTest(t *testing.T, ref reference.Named) (*http.Client, *httptest.Server) {
	t.Helper()

	server := vtesting.NewTestDockerServer(t)

	auth := vtesting.NewAuth(server)

	client, err := auth.ToClient(context.TODO(), ref)
	require.NoError(t, err, "failed to create client: %s", err)

	vtesting.UpdateClient(client, server)

	return client, server
}
