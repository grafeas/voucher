package diy

import (
	"context"
	"time"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/docker"
)

// check is a check that verifies if the passed image was built
// by us.
type check struct {
	auth voucher.Auth
}

// SetAuth sets the authentication system that this check will use
// for its run.
func (d *check) SetAuth(auth voucher.Auth) {
	d.auth = auth
}

// check checks if an image was built by a trusted source
func (d *check) Check(i voucher.ImageData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

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
