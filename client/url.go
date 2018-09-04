package client

import (
	"net/url"
	"path"
)

// toVoucherURL adds the check to the URL properly and returns
// a string containing the full URL.
func toVoucherURL(voucherURL *url.URL, checkname string) string {
	if nil == voucherURL {
		return "/" + checkname
	}

	voucherURL.Path = path.Join(voucherURL.Path, checkname)
	return voucherURL.String()
}
