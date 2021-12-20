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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	httptransport "google.golang.org/api/transport/http"
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

// NewClientContext creates a new Client set to connect to the passed
// hostname.
func NewClientContext(ctx context.Context, voucherURL string, options ...Option) (*Client, error) {
	if voucherURL == "" {
		return nil, errNoHost
	}

	u, err := url.Parse(voucherURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse voucher hostname: %w", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	client := &Client{
		url:        u,
		httpClient: &http.Client{},
		userAgent:  DefaultUserAgent,
	}
	for _, opt := range options {
		if err := opt(ctx, client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

type Option func(context.Context, *Client) error

// WithHTTPClient customizes the http.Client used by the client.
// Customize the client's Transport to do arbitrary authentication.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(_ context.Context, c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithBasicAuth sets the username and password to use for the client.
func WithBasicAuth(username, password string) Option {
	return func(_ context.Context, c *Client) error {
		c.username = username
		c.password = password
		return nil
	}
}

// WithIDTokenAuth configures the client to use an ID token.
func WithIDTokenAuth() Option {
	return func(ctx context.Context, c *Client) error {
		idClient, err := idtoken.NewClient(ctx, c.url.String())
		if err != nil {
			return err
		}
		c.httpClient = idClient
		return nil
	}
}

// WithDefaultTokenAuth configures the client to use Google's default token.
func WithDefaultTokenAuth() Option {
	return func(ctx context.Context, c *Client) error {
		src, err := google.DefaultTokenSource(ctx)
		if err != nil {
			return fmt.Errorf("error getting default token source: %w", err)
		}
		ts := oauth2.ReuseTokenSource(nil, &idTokenSource{TokenSource: src, audience: c.url.String()})

		transport, err := httptransport.NewTransport(ctx, http.DefaultTransport, option.WithTokenSource(ts))
		if err != nil {
			return fmt.Errorf("error creating client: %w", err)
		}
		c.httpClient = &http.Client{Transport: transport}
		return nil
	}
}

// WithUserAgent sets the User-Agent header for the client.
func WithUserAgent(userAgent string) Option {
	return func(_ context.Context, c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}

// NewAuthClient creates a new auth Client set to connect to the passed
// hostname using tokens.
// Deprecated: use NewClientContext and the WithIDTokenAuth option instead
func NewAuthClient(voucherURL string) (*Client, error) {
	return NewClientContext(context.Background(), voucherURL, WithIDTokenAuth())
}

// NewClient creates a new auth Client.
// Deprecated: use NewClientContext
func NewClient(voucherURL string) (*Client, error) {
	return NewClientContext(context.Background(), voucherURL)
}

// SetBasicAuth adds the username and password to the Client struct
// Deprecated: use the WithBasicAuth option instead
func (c *Client) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
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
