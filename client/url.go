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

	// Copy our URL, so we are not modifying the original.
	newVoucherURL := (*voucherURL)

	newVoucherURL.Path = path.Join(newVoucherURL.Path, checkname)
	return newVoucherURL.String()
}
