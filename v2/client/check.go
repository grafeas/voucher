package client

import (
	"context"
	"net/url"
	"path"

	"github.com/docker/distribution/reference"
	voucher "github.com/grafeas/voucher/v2"
)

// Check executes a request to a Voucher server, to the appropriate check URI, and
// with the passed reference.Canonical. Returns a voucher.Response and an error.
func (c *Client) Check(ctx context.Context, check string, image reference.Canonical) (voucher.Response, error) {
	url := toVoucherCheckURL(c.url, check)
	resp, err := c.doVoucherRequest(ctx, url, image)
	if err != nil {
		return voucher.Response{}, err
	}
	return *resp, nil
}

// toVoucherCheckURL adds the check to the URL properly and returns
// a string containing the full URL.
func toVoucherCheckURL(voucherURL *url.URL, checkname string) string {
	if nil == voucherURL {
		return "/" + checkname
	}

	// Copy our URL, so we are not modifying the original.
	newVoucherURL := (*voucherURL)

	newVoucherURL.Path = path.Join(newVoucherURL.Path, checkname)
	return newVoucherURL.String()
}
