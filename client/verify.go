package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher"
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
