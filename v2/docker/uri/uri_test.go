package uri

import (
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testHostname    = "gcr.io"
	testProject     = "test/project"
	testDigest      = "sha256:cb749360c5198a55859a7f335de3cf4e2f64b60886a2098684a2f9c7ffca81f2"
	testBlobURL     = "https://" + testHostname + "/v2/" + testProject + "/blobs/" + testDigest
	testManifestURL = "https://" + testHostname + "/v2/" + testProject + "/manifests/" + testDigest
	testTokenURL    = "https://" + testHostname + "/v2/token?scope=repository%3Atest%2Fproject%3A%2A&service=gcr.io"
)

func TestGetBaseURI(t *testing.T) {
	named, err := reference.ParseNamed(testHostname + "/" + testProject + "@" + testDigest)
	require.NoError(t, err, "failed to parse uri: %s", err)

	assert.Equal(t, testTokenURL, GetTokenURI(named))

	canonicalRef, ok := named.(reference.Canonical)
	require.True(t, ok)

	assert.Equal(t, string(canonicalRef.Digest()), testDigest)
	hostname, path := reference.SplitHostname(canonicalRef)
	assert.Equal(t, hostname, "gcr.io")
	assert.Equal(t, path, testProject)
	assert.Equal(t, testBlobURL, GetBlobURI(canonicalRef, canonicalRef.Digest()))
	assert.Equal(t, testManifestURL, GetDigestManifestURI(canonicalRef))
}
