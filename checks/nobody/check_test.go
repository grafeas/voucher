package nobody

import (
	"context"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafeas/voucher"
	vtesting "github.com/grafeas/voucher/testing"
)

func TestNobodyCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	auth := vtesting.NewAuth(server)

	nobodyCheck := new(check)
	nobodyCheck.SetAuth(auth)

	i := vtesting.NewTestReference(t)

	pass, err := nobodyCheck.Check(context.Background(), i)

	require.NoErrorf(t, err, "check failed with error: %s", err)
	assert.True(t, pass, "check failed when it should have passed")
}

func TestThirdPartyImage(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	auth := vtesting.NewAuth(server)

	name := "registry-1.docker.io/library/alpine@sha256:4e01ddea8def856ba9fee17668fa0b2e45a8bc78127b7ab6cf921f6d6fd86ac9"
	ref, err := reference.Parse(name)
	require.NoErrorf(t, err, "could not make image reference (\"%s\"): %s", name, err)

	refCanonical, ok := ref.(reference.Canonical)
	require.True(t, ok, "could not convert reference to Canonical reference")

	nobodyCheck := new(check)
	nobodyCheck.SetAuth(auth)

	pass, err := nobodyCheck.Check(context.Background(), refCanonical)

	require.Error(t, err, "check should have failed with error, but didn't")
	assert.Contains(t, err.Error(), "auth failed: does not match domain", "check should have failed due to invalid image reference domain, but didn't")
	assert.False(t, pass, "check passed when it should have failed")
}

func TestNobodyCheckWithNoAuth(t *testing.T) {
	i := vtesting.NewTestReference(t)

	nobodyCheck := new(check)

	// run check without setting up Auth.
	pass, err := nobodyCheck.Check(context.Background(), i)
	require.Equal(t, err, voucher.ErrNoAuth, "check should have failed due to lack of Auth, but didn't")
	assert.False(t, pass, "check passed when it should have failed due to no Auth")
}

func TestFailingNobodyCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	auth := vtesting.NewAuth(server)

	i := vtesting.NewNobodyBadTestReference(t)

	nobodyCheck := new(check)
	nobodyCheck.SetAuth(auth)

	pass, err := nobodyCheck.Check(context.Background(), i)

	require.NoError(t, err, "check should have failed with error, but didn't")
	assert.False(t, pass, "check passed when it should have failed")
}
