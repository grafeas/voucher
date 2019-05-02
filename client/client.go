package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Shopify/voucher"
	"github.com/docker/distribution/reference"
)

var errNoHost = errors.New("cannot create client with empty hostname")

// VoucherClient is a client for the Voucher API.
type VoucherClient struct {
	Hostname   *url.URL
	httpClient *http.Client
	username   string
	password   string
}

// Check executes a request to a Voucher server, to the appropriate check URI, and
// with the passed reference.Canonical. Returns a voucher.Response and an error.
func (c *VoucherClient) Check(check string, image reference.Canonical) (voucher.Response, error) {
	var checkResp voucher.Response
	var buffer bytes.Buffer

	err := json.NewEncoder(&buffer).Encode(voucher.Request{
		ImageURL: image.String(),
	})
	if err != nil {
		return checkResp, fmt.Errorf("could not parse image, error: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, toVoucherURL(c.Hostname, check), &buffer)
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
		b, err := ioutil.ReadAll(resp.Body)
		if nil == err {
			err = fmt.Errorf("failed to get response: %s", strings.TrimSpace(string(b)))
		}
		return checkResp, err
	}

	err = json.NewDecoder(resp.Body).Decode(&checkResp)
	return checkResp, err
}

// NewClient creates a new VoucherClient set to connect to the passed
// hostname, and with the passed timeout.
func NewClient(hostname string, timeout time.Duration) (*VoucherClient, error) {
	var err error

	if "" == hostname {
		return nil, errNoHost
	}

	client := new(VoucherClient)
	client.httpClient = &http.Client{
		Timeout: timeout,
	}

	client.Hostname, err = url.Parse(hostname)
	if nil != err {
		return nil, fmt.Errorf("could not parse voucher hostname: %s", err)
	}

	if "" == client.Hostname.Scheme {
		client.Hostname.Scheme = "https"
	}

	return client, nil
}

// SetBasicAuth adds the username and password to the VoucherClient struct
func (c *VoucherClient) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
}
