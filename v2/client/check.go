package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/docker/distribution/reference"
	voucher "github.com/grafeas/voucher/v2"
)

// Check executes a request to a Voucher server, to the appropriate check URI, and
// with the passed reference.Canonical. Returns a voucher.Response and an error.
func (c *Client) Check(ctx context.Context, check string, image reference.Canonical) (voucher.Response, error) {
	var checkResp voucher.Response
	var buffer bytes.Buffer

	err := json.NewEncoder(&buffer).Encode(voucher.Request{
		ImageURL: image.String(),
	})
	if err != nil {
		return checkResp, fmt.Errorf("could not parse image, error: %s", err)
	}

	var req *http.Request

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		toVoucherCheckURL(c.hostname, check),
		&buffer,
	)
	if nil != err {
		return checkResp, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	resp, err := c.httpClient.Do(req)
	if nil != err {
		return checkResp, err
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if nil == err {
			err = fmt.Errorf("failed to get response: %s", strings.TrimSpace(string(b)))
		}
		return checkResp, err
	}

	err = json.NewDecoder(resp.Body).Decode(&checkResp)
	return checkResp, err
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
