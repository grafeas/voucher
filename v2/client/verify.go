package client

import (
	"context"
	"net/url"
	"path"

	"github.com/docker/distribution/reference"

	voucher "github.com/grafeas/voucher/v2"
)

func (c *Client) Verify(ctx context.Context, check string, image reference.Canonical) (voucher.Response, error) {
	url := toVoucherVerifyURL(c.url, check)
	resp, err := c.doVoucherRequest(ctx, url, image)
	if err != nil {
		return voucher.Response{}, err
	}
	return *resp, nil
}

// toVoucherVerifyURL adds the check and verify API path to the URL properly
// and returns a string containing the full URL.
func toVoucherVerifyURL(voucherURL *url.URL, checkname string) string {
	if nil == voucherURL {
		return "/" + checkname + "/verify"
	}

	// Copy our URL, so we are not modifying the original.
	newVoucherURL := (*voucherURL)

	newVoucherURL.Path = path.Join(newVoucherURL.Path, checkname, "verify")
	return newVoucherURL.String()
}
