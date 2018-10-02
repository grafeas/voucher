package diy

import (
	"testing"

	"github.com/Shopify/voucher"
	vtesting "github.com/Shopify/voucher/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDIYCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	i := vtesting.NewTestReference(t)

	diyCheck := new(check)
	diyCheck.SetAuth(vtesting.NewAuth(server))

	pass, err := diyCheck.Check(i)

	assert.NoErrorf(t, err, "check failed with error: %s", err)
	assert.True(t, pass, "check failed when it should have passed")
}

func TestDIYCheckWithNoAuth(t *testing.T) {
	i := vtesting.NewTestReference(t)

	diyCheck := new(check)

	// run check without setting up Auth.
	pass, err := diyCheck.Check(i)

	assert.Equal(t, err, voucher.ErrNoAuth, "check should have failed due to lack of Auth, but didn't")
	assert.False(t, pass, "check passed when it should have failed due to no Auth")
}

func TestFailingDIYCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	auth := vtesting.NewAuth(server)

	diyCheck := new(check)
	diyCheck.SetAuth(auth)

	i := vtesting.NewBadTestReference(t)

	pass, err := diyCheck.Check(i)

	require.Error(t, err, "check should have failed with error, but didn't")
	assert.Containsf(t, err.Error(), "image doesn't exist", "check error format is incorrect: \"%s\"", err)
	assert.False(t, pass, "check passed when it should have failed")
}
