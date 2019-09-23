package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	_, err := NewClient("")
	assert.Equalf(err, errNoHost, "should have been a no-host error, is actually: %s", err)

	client, err := NewClient("localhost")
	assert.NoErrorf(err, "failed to create client: %s", err)
	assert.NotNilf(client.hostname, "client hostname URL is nil")
	assert.Equal(client.hostname.String(), "https://localhost")

	_, err = NewClient(":localhost")
	require.Contains(t, err.Error(), "could not parse voucher hostname", "failed to create client: ", err)
}

func TestVoucherURL(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("localhost")
	require.NoError(t, err, "failed to create client: ", err)

	allTestURL := toVoucherURL(client.hostname, "all")
	assert.Equalf(allTestURL, "https://localhost/all", "url is incorrect, should be \"%s\" instead of \"%s\"", "https://localhost/all", allTestURL)

	allEmptyURL := toVoucherURL(nil, "all")
	assert.Equal(allEmptyURL, "/all")
}

func TestVoucherBasicAuth(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("localhost")
	assert.NoErrorf(err, "failed to create client: %s", err)
	assert.Equal(client.username, "", "username already set in client")
	assert.Equal(client.password, "", "password already set in client")

	client.SetBasicAuth("username", "password")

	assert.Equal(client.username, "username", "username incorrect in client: %s", client.username)
	assert.Equal(client.password, "password", "password incorrect in client: %s", client.password)
}
