package client

import (
	"net/url"
	"path"
)

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
