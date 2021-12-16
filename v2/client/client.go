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
	voucher "github.com/grafeas/voucher/v2"
	"google.golang.org/api/idtoken"
)

var errNoHost = errors.New("cannot create client with empty hostname")

// Client is a client for the Voucher API.
type Client struct {
	url        *url.URL
	httpClient *http.Client
	username   string
	password   string
	userAgent  string
}

const DefaultUserAgent = "voucher-client/2"

// NewClient creates a new Client set to connect to the passed
// hostname.
func NewClient(voucherURL string) (*Client, error) {
	return newClient(voucherURL, &http.Client{})
}

// NewAuthClient creates a new auth Client set to connect to the passed
// hostname using tokens.
func NewAuthClient(voucherURL string) (*Client, error) {
	authClient, err := idtoken.NewClient(context.Background(), voucherURL)
	if err != nil {
		return nil, err
	}
	return newClient(voucherURL, authClient)
}

// SetBasicAuth adds the username and password to the Client struct
func (c *Client) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
}

// SetUserAgent customizes the user agent used by the client
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// CopyURL returns a copy of this client's URL
func (c *Client) CopyURL() *url.URL {
	urlCopy := (*c.url)
	return &urlCopy
}

func (c *Client) newVoucherRequest(ctx context.Context, url string, image reference.Canonical) (*http.Request, error) {
	voucherReq := voucher.Request{
		ImageURL: image.String(),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(voucherReq); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	return req, nil
}

func (c *Client) doVoucherRequest(ctx context.Context, url string, image reference.Canonical) (*voucher.Response, error) {
	req, err := c.newVoucherRequest(ctx, url, image)
	if err != nil {
		return nil, fmt.Errorf("could create voucher request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		b, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err = fmt.Errorf("failed to get response: %s", strings.TrimSpace(string(b)))
		}
		return nil, err
	}

	var voucherResp voucher.Response
	if err := json.NewDecoder(resp.Body).Decode(&voucherResp); err != nil {
		return nil, err
	}
	return &voucherResp, nil
}

func newClient(voucherURL string, httpClient *http.Client) (*Client, error) {
	if voucherURL == "" {
		return nil, errNoHost
	}

	u, err := url.Parse(voucherURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse voucher hostname: %s", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	client := &Client{
		url:        u,
		httpClient: httpClient,
		userAgent:  DefaultUserAgent,
	}
	return client, nil
}
