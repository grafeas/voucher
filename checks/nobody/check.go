package nobody

import (
	"context"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/docker"
)

// check is for verifying that the passed image does not run as
// root or user 0.
type check struct {
	auth voucher.Auth
}

// SetAuth sets the authentication system that this check will use
// for its run.
func (n *check) SetAuth(auth voucher.Auth) {
	n.auth = auth
}

// Check verifies if the image runs as root and returns a boolean (true if
// the user is not root, false otherwise) and an error as response.
func (n *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	if nil == n.auth {
		return false, voucher.ErrNoAuth
	}

	client, err := n.auth.ToClient(ctx, i)
	if nil != err {
		return false, err
	}

	imageConfig, err := docker.RequestImageConfig(client, i)

	if nil != err {
		return false, err
	}

	return !imageConfig.RunsAsRoot(), nil
}

func init() {
	voucher.RegisterCheckFactory("nobody", func() voucher.Check {
		return new(check)
	})
}
