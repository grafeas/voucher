package client

import (
	"time"

	"github.com/Shopify/voucher"
	"github.com/docker/distribution/reference"
)

const timeout = 120 * time.Second

// SignImage takes an image URL (which includes the registry and the image
// digest) as well as the check to be performed and makes a call to Voucher
// to run the specified checks. It returns an error, or nil if no errors
// occur and the check was successful.
func SignImage(hostname string, image reference.Canonical, check string) (voucher.Response, error) {
	client, err := NewClient(hostname, timeout)
	if nil != err {
		return voucher.Response{}, err
	}

	checkResp, err := client.Check(check, image)
	if nil != err {
		return voucher.Response{}, err
	}
	return checkResp, nil
}
