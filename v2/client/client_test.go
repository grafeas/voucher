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
			c, err := client.NewClient(tc.input)
			if tc.errMsg != "" {
				assert.Equal(t, tc.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.hostname, c.URL().String())
			}
		})
	}
}

const image = "gcr.io/project/image@sha256:0000000000000000000000000000000000000000000000000000000000000000"

func TestVoucher_Check(t *testing.T) {
	v := &mockVoucher{}
	v.checks = append(v.checks, &voucher.Response{Image: image, Success: true})
	srv := httptest.NewServer(v)
	defer srv.Close()

	c, err := client.NewClient(srv.URL)
	require.NoError(t, err)
	res, err := c.Check(context.Background(), "diy", canonical(t, image))
	require.NoError(t, err)
	assert.True(t, res.Success)
}

func TestVoucher_Verify(t *testing.T) {
	v := &mockVoucher{}
	v.verifications = append(v.verifications, &voucher.Response{Image: image, Success: true})
	srv := httptest.NewServer(v)
	defer srv.Close()

	c, err := client.NewClient(srv.URL)
	require.NoError(t, err)
	res, err := c.Verify(context.Background(), "diy", canonical(t, image))
	require.NoError(t, err)
	assert.True(t, res.Success)
}

type mockVoucher struct {
	checks        []*voucher.Response
	verifications []*voucher.Response
}

func (v *mockVoucher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
