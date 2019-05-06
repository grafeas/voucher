package docker

import (
	"testing"

	vtesting "github.com/Shopify/voucher/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestManifest(t *testing.T) {
	ref := vtesting.NewTestReference(t)

	client, server := PrepareDockerTest(t, ref)
	defer server.Close()

	manifest, err := RequestManifest(client, ref)
	require.NoError(t, err, "failed to get manifest: %s", err)

	assert.Equal(t, vtesting.NewTestManifest(), manifest)
}

func TestRequestBadManifest(t *testing.T) {
	ref := vtesting.NewBadTestReference(t)

	client, server := PrepareDockerTest(t, ref)
	defer server.Close()

	_, err := RequestManifest(client, ref)
	assert.NotNilf(t, err, "should have failed to get manifest, but didn't")
	assert.Contains(t, err.Error(), "failed to load resource with status \"404 Not Found\":")
}

func TestRateLimitedBadManifest(t *testing.T) {
	ref := vtesting.NewRateLimitedTestReference(t)

	client, server := PrepareDockerTest(t, ref)
	defer server.Close()

	_, err := RequestManifest(client, ref)
	assert.NotNilf(t, err, "should have failed to get manifest, but didn't")
	assert.Contains(t, err.Error(), "failed to load resource with status \"200 OK\": "+vtesting.RateLimitOutput)
}
