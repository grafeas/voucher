package sbom

import (
	"testing"

	"github.com/docker/distribution/reference"
	sbomgcr "github.com/grafeas/voucher/v2/sbomgcr"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasSBOM(t *testing.T) {
	// TODO:CS come back and actually make the mocks for these, don't merge this in! It's making network requests!
	mockService := sbomgcr.NewGCRService()
	mockSBOMClient := sbomgcr.NewClient(mockService)
	mockCheck := check{sbomClient: mockSBOMClient}

	img, digest := "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui@sha256:551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f", "sha256:551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f"
	ref := getCanonicalRef(t, img, digest)

	sbom := mockCheck.hasSBOM(ref)

	assert.True(t, sbom)
}

func getCanonicalRef(t *testing.T, img string, digestStr string) reference.Canonical {
	named, err := reference.ParseNamed(img)
	require.NoError(t, err, "named")
	canonicalRef, err := reference.WithDigest(named, digest.Digest(digestStr))
	require.NoError(t, err, "canonicalRef")
	return canonicalRef
}
