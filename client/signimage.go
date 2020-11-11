package client

import (
	"context"
	"time"

	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher"
)

const timeout = 120 * time.Second

// SignImage takes an image URL (which includes the registry and the image
// digest) as well as the check to be performed and makes a call to Voucher
// to run the specified checks. It returns an error, or nil if no errors
// occur and the check was successful.
func SignImage(hostname string, image reference.Canonical, check string) (voucher.Response, error) {
	client, err := NewClient(hostname)
	if nil != err {
		return voucher.Response{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	checkResp, err := client.Check(ctx, check, image)
	if nil != err {
		return voucher.Response{}, err
	}
	return checkResp, nil
}
