package docker

import (
	"testing"

	"github.com/docker/distribution/reference"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vtesting "github.com/grafeas/voucher/testing"
)

func TestGetDigestFromTagged(t *testing.T) {
	ref := vtesting.NewTestReference(t)

	taggedRef, err := reference.WithTag(ref, "latest")
	require.NoErrorf(t, err, "failed to get tagged reference: %s", err)

	client, server := vtesting.PrepareDockerTest(t, taggedRef)
	defer server.Close()

	imageDigest, err := GetDigestFromTagged(client, taggedRef)
	require.NoErrorf(t, err, "failed to get digest reference: %s", err)

	assert.Equal(t, digest.Digest("sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da"), imageDigest)
}

func TestGetBadDigestFromTagged(t *testing.T) {
	ref := vtesting.NewBadTestReference(t)

	taggedRef, err := reference.WithTag(ref, "latest")
	require.NoErrorf(t, err, "failed to get tagged reference: %s", err)

	client, server := vtesting.PrepareDockerTest(t, taggedRef)
	defer server.Close()

	imageDigest, err := GetDigestFromTagged(client, taggedRef)
	assert.NotNilf(t, err, "should have failed to get digest, but didn't")
	assert.Equal(t, digest.Digest(""), imageDigest)
	assert.Contains(t, err.Error(), "failed to load resource with status \"404 Not Found\":")
}

func TestRequestV1Digest(t *testing.T) {
	ref := vtesting.NewTestSchema1SignedReference(t)

	taggedRef, err := reference.WithTag(ref, "latest")
	require.NoErrorf(t, err, "failed to get tagged reference: %s", err)

	client, server := vtesting.PrepareDockerTest(t, taggedRef)
	defer server.Close()

	imageDigest, err := GetDigestFromTagged(client, taggedRef)
	require.NoErrorf(t, err, "failed to get digest reference: %s", err)

	assert.Equal(t, digest.Digest("sha256:18e6e7971438ab792d13563dcd8972acf4445bc0dcfdff84a6374d63a9c3ed62"), imageDigest)
}
