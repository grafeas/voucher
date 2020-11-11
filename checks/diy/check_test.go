package diy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafeas/voucher"
	vtesting "github.com/grafeas/voucher/testing"
)

func TestDIYCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	i := vtesting.NewTestReference(t)

	diyCheck := new(check)
	diyCheck.SetAuth(vtesting.NewAuth(server))
	diyCheck.SetValidRepos([]string{
		i.Name(),
	})

	pass, err := diyCheck.Check(context.Background(), i)

	assert.NoErrorf(t, err, "check failed with error: %s", err)
	assert.True(t, pass, "check failed when it should have passed")
}

func TestDIYCheckWithInvalidRepo(t *testing.T) {
	i := vtesting.NewTestReference(t)

	diyCheck := new(check)

	// run check without setting up valid repos.
	pass, err := diyCheck.Check(context.Background(), i)

	assert.Equal(t, err, ErrNotFromRepo, "check should have failed due to image not being from a valid repo, but didn't")
	assert.False(t, pass, "check passed when it should have failed due to image being from an invalid repo")
}

func TestDIYCheckWithNoAuth(t *testing.T) {
	i := vtesting.NewTestReference(t)

	diyCheck := new(check)
	diyCheck.SetValidRepos([]string{
		i.Name(),
	})

	// run check without setting up Auth.
	pass, err := diyCheck.Check(context.Background(), i)

	assert.Equal(t, err, voucher.ErrNoAuth, "check should have failed due to lack of Auth, but didn't")
	assert.False(t, pass, "check passed when it should have failed due to no Auth")
}

func TestFailingDIYCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	auth := vtesting.NewAuth(server)

	i := vtesting.NewBadTestReference(t)

	diyCheck := new(check)
	diyCheck.SetAuth(auth)
	diyCheck.SetValidRepos([]string{
		i.Name(),
	})

	pass, err := diyCheck.Check(context.Background(), i)

	require.Error(t, err, "check should have failed with error, but didn't")
	assert.Containsf(t, err.Error(), "image doesn't exist", "check error format is incorrect, should be \"image doesn't exist\": \"%s\"", err)
	assert.False(t, pass, "check passed when it should have failed")
}
