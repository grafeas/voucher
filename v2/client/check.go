package client

import (
	"context"
	"path"

	voucher "github.com/grafeas/voucher/v2"
)

// Check executes a request to a Voucher server, to the appropriate check URI, and
// with the passed image. Returns a voucher.Response and an error.
func (c *Client) Check(ctx context.Context, check string, image string) (voucher.Response, error) {
	url := c.toVoucherCheckURL(check)
	resp, err := c.doVoucherRequest(ctx, url, image)
	if err != nil {
		return voucher.Response{}, err
	}
	return *resp, nil
}

func (c *Client) toVoucherCheckURL(checkname string) string {
	newVoucherURL := c.CopyURL()
	newVoucherURL.Path = path.Join(newVoucherURL.Path, checkname)
	return newVoucherURL.String()
}
