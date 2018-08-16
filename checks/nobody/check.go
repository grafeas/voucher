package nobody

import (
	voucher "github.com/Shopify/voucher"
	"github.com/Shopify/voucher/docker"
)

// check is for verifying that the passed image does not run as
// root or user 0.
type check struct {
}

// Check verifies if the image runs as root and returns a boolean (true if
// the user is not root, false otherwise) and an error as response.
func (n *check) Check(i voucher.ImageData) (bool, error) {
	gcloudToken, err := voucher.GetAccessToken()
	if nil != err {
		return false, err
	}

	oauthToken, err := docker.Auth(gcloudToken, i)
	if nil != err {
		return false, err
	}

	imageConfig, err := docker.RequestImageConfig(oauthToken, i)
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
