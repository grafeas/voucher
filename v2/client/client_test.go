package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cases := map[string]struct {
		input    string
		hostname string
		errMsg   string
	}{
		"no host": {
			input:  "",
			errMsg: "cannot create client with empty hostname",
		},
		"no scheme": {
			input:    "localhost",
			hostname: "https://localhost",
		},
		"bad url": {
			input:  ":localhost",
			errMsg: `could not parse voucher hostname: parse ":localhost": missing protocol scheme`,
		},
	}
	for label, tc := range cases {
		t.Run(label, func(t *testing.T) {
			c, err := client.NewClient(context.Background(), tc.input)
			if tc.errMsg != "" {
				assert.EqualError(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.hostname, c.CopyURL().String())
			}
		})
	}
}

const image = "gcr.io/project/image@sha256:0000000000000000000000000000000000000000000000000000000000000000"

func TestVoucher_Check(t *testing.T) {
	v := &mockVoucher{t: t}
	v.checks = append(v.checks, &voucher.Response{Image: image, Success: true})
	srv := httptest.NewServer(v)
	defer srv.Close()

	c, err := client.NewClient(context.Background(), srv.URL)
	require.NoError(t, err)
	res, err := c.Check(context.Background(), "diy", canonical(t, image))
	require.NoError(t, err)
	assert.True(t, res.Success)
}

func TestVoucher_CustomUserAgent(t *testing.T) {
	const customUserAgent = "my-awesome-voucher/1.0"
	v := &mockVoucher{t: t}
	v.ua = customUserAgent
	v.checks = append(v.checks, &voucher.Response{Image: image, Success: true})
	srv := httptest.NewServer(v)
	defer srv.Close()

	c, err := client.NewClient(context.Background(), srv.URL, client.WithUserAgent(customUserAgent))
	require.NoError(t, err)
	res, err := c.Check(context.Background(), "diy", canonical(t, image))
	require.NoError(t, err)
	assert.True(t, res.Success)
}

func TestVoucher_Verify(t *testing.T) {
	v := &mockVoucher{t: t}
	v.verifications = append(v.verifications, &voucher.Response{Image: image, Success: true})
	srv := httptest.NewServer(v)
	defer srv.Close()

	c, err := client.NewClient(context.Background(), srv.URL)
	require.NoError(t, err)
	res, err := c.Verify(context.Background(), "diy", canonical(t, image))
	require.NoError(t, err)
	assert.True(t, res.Success)
}

type mockVoucher struct {
	t             *testing.T
	ua            string
	checks        []*voucher.Response
	verifications []*voucher.Response
}

func (v *mockVoucher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hdr := r.Header
	assert.Equal(v.t, "application/json", hdr.Get("Content-Type"))
	if v.ua != "" {
		assert.Equal(v.t, v.ua, hdr.Get("User-Agent"))
	} else {
		assert.Equal(v.t, client.DefaultUserAgent, hdr.Get("User-Agent"))
	}

	var req voucher.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var search []*voucher.Response
	if strings.Contains(r.URL.Path, "/verify") {
		search = v.verifications
	} else {
		search = v.checks
	}

	var res *voucher.Response
	for _, r := range search {
		if r.Image == req.ImageURL {
			res = r
			break
		}
	}
	if res == nil {
		res = &voucher.Response{
			Image:   req.ImageURL,
			Success: false,
		}
	}
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func canonical(t *testing.T, image string) reference.Canonical {
	ref, err := reference.Parse(image)
	require.NoError(t, err)
	canonical, ok := ref.(reference.Canonical)
	require.True(t, ok)
	return canonical
}
