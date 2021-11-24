package docker

import (
	"context"
	"os"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

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
	// FIXME: follow the stubbed hub style of the abov ; oreilly-well-do-it-live.jpg
	token := &oauth2.Token{AccessToken: os.Getenv("REGISTRY_TOKEN")}
	c := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	ref, err := reference.Parse("registry.k8s.pwagner.net/dockerhub/library/ubuntu@sha256:cc8f713078bfddfe9ace41e29eb73298f52b2c958ccacd1b376b9378e20906ef")
	require.NoError(t, err)

	ic, err := RequestImageConfig(c, ref.(reference.Canonical))
	require.NoError(t, err)
	assert.True(t, ic.RunsAsRoot())
}
