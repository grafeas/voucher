package client

import (
	"context"
	"path"

	"github.com/docker/distribution/reference"

	voucher "github.com/grafeas/voucher/v2"
)

func (c *Client) Verify(ctx context.Context, check string, image reference.Canonical) (voucher.Response, error) {
	url := c.toVoucherVerifyURL(check)
	resp, err := c.doVoucherRequest(ctx, url, image)
	if err != nil {
		return voucher.Response{}, err
	}
	return *resp, nil
}

func (c *Client) toVoucherVerifyURL(checkname string) string {
	newVoucherURL := c.URL()
	newVoucherURL.Path = path.Join(newVoucherURL.Path, checkname, "verify")
	return newVoucherURL.String()
}
