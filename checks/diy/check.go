package diy

import (
	"context"
	"errors"
	"strings"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/docker"
)

// ErrNotFromRepo is returned when an image does not match one of the valid
// repo paths.
var ErrNotFromRepo = errors.New("image is not from a valid repo")

// check is a check that verifies if the passed image was built
// by us.
type check struct {
	auth       voucher.Auth
	validRepos []string
}

// SetValidRepos sets the repos that images must be in to get signed by the
// DIY check.
func (d *check) SetValidRepos(repos []string) {
	d.validRepos = repos
}

// SetAuth sets the authentication system that this check will use
// for its run.
func (d *check) SetAuth(auth voucher.Auth) {
	d.auth = auth
}

// isFromValidrepo returns true if the passed image is from a valid repo.
func (d *check) isFromValidRepo(i voucher.ImageData) bool {
	for _, repo := range d.validRepos {
		if strings.HasPrefix(i.Name(), repo) {
			return true
		}
	}
	return false
}

// check checks if an image was built by a trusted source
func (d *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	if !d.isFromValidRepo(i) {
		return false, ErrNotFromRepo
	}

	if nil == d.auth {
		return false, voucher.ErrNoAuth
	}

	client, err := d.auth.ToClient(ctx, i)
	if nil != err {
		return false, err
	}

	_, err = docker.RequestImageConfig(client, i)
	if nil != err {
		return false, err
	}

	return true, nil
}

func init() {
	voucher.RegisterCheckFactory("diy", func() voucher.Check {
		return new(check)
	})
}
