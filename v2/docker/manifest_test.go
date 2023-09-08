package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafeas/voucher/v2/docker/schema2"
	vtesting "github.com/grafeas/voucher/v2/testing"
)

func TestRequestManifest(t *testing.T) {
	ref := vtesting.NewTestReference(t)

	client, server := vtesting.PrepareDockerTest(t, ref)
	defer server.Close()

	manifest, err := RequestManifest(client, ref)
	require.NoError(t, err)

	schema2Manifest, err := schema2.ToManifest(client, ref, manifest)
	require.NoError(t, err)

	assert.Equal(
		t,
		vtesting.NewTestManifest().Manifest,
		schema2Manifest,
	)
}

func TestRequestBadManifest(t *testing.T) {
	ref := vtesting.NewBadTestReference(t)

	client, server := vtesting.PrepareDockerTest(t, ref)
	defer server.Close()

	_, err := RequestManifest(client, ref)
	require.NotNilf(t, err, "should have failed to get manifest, but didn't")
	assert.Equal(t,
		NewManifestErrorWithRequest("404 Not Found", []byte("image doesn't exist\n")),
		err,
	)
}

func TestRateLimitedBadManifest(t *testing.T) {
	ref := vtesting.NewRateLimitedTestReference(t)

	client, server := vtesting.PrepareDockerTest(t, ref)
	defer server.Close()

	_, err := RequestManifest(client, ref)
	assert.NotNilf(t, err, "should have failed to get manifest, but didn't")
	assert.Equal(t,
		NewManifestErrorWithRequest("200 OK", []byte(vtesting.RateLimitOutput+"\n")),
		err,
	)
}

func TestRequestManifestList(t *testing.T) {
	ref := vtesting.NewTestManifestListReference(t)

	client, server := vtesting.PrepareDockerTest(t, ref)
	defer server.Close()

	manifest, err := RequestManifest(client, ref)
	require.NoError(t, err)

	schema2Manifest, err := schema2.ToManifest(client, ref, manifest)
	require.NoError(t, err)

	assert.Equal(
		t,
		vtesting.NewTestManifest().Manifest,
		schema2Manifest,
	)
}
