package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher"
)

var errNoHost = errors.New("cannot create client with empty hostname")

// Client is a client for the Voucher API.
type Client struct {
	hostname   *url.URL
	httpClient *http.Client
	username   string
	password   string
}

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

// SetBasicAuth adds the username and password to the Client struct
func (c *Client) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
}

// NewClient creates a new Client set to connect to the passed
// hostname.
func NewClient(hostname string) (*Client, error) {
	var err error

	if "" == hostname {
		return nil, errNoHost
	}

	hostnameURL, err := url.Parse(hostname)
	if nil != err {
		return nil, fmt.Errorf("could not parse voucher hostname: %s", err)
	}

	if "" == hostnameURL.Scheme {
		hostnameURL.Scheme = "https"
	}

	client := &Client{
		hostname:   hostnameURL,
		httpClient: &http.Client{},
	}

	return client, nil
}
