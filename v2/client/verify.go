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

func (c *Client) Verify(ctx context.Context, check string, image reference.Canonical) (voucher.Response, error) {
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
		toVoucherVerifyURL(c.hostname, check),
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
