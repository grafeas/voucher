package diy

import (
	voucher "github.com/Shopify/voucher"
	"github.com/Shopify/voucher/docker"
)

// check is a check that verifies if the passed image was built
// by us.
type check struct {
}

// check checks if an image was built by a trusted source
func (d *check) Check(i voucher.ImageData) (bool, error) {
	gcloudToken, err := voucher.GetAccessToken()
	if nil != err {
		return false, err
	}

	oauthToken, err := docker.Auth(gcloudToken, i)
	if nil != err {
		return false, err
	}

	_, err = docker.RequestImageConfig(oauthToken, i)
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
