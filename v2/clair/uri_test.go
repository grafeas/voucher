package clair

import (
	"testing"

	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

const (
	testHostname                = "clair.example.com"
	testHostnameWithProtocol    = "http://clair.example.com"
	testDigest                  = "sha256:cb749360c5198a55859a7f335de3cf4e2f64b60886a2098684a2f9c7ffca81f2"
	testNewLayerURI             = "https://" + testHostname + "/v1/layers"
	testGetLayerURI             = "https://" + testHostname + "/v1/layers/" + testDigest + "?vulnerabilities"
	testNewLayerURIWithProtocol = testHostnameWithProtocol + "/v1/layers"
)

func TestGetURIs(t *testing.T) {
	assert.Equal(t, testNewLayerURI, GetNewLayerURI(testHostname))
	assert.Equal(t, testGetLayerURI, GetLayerURI(testHostname, digest.Digest(testDigest)))
	assert.Equal(t, testNewLayerURIWithProtocol, GetNewLayerURI(testHostnameWithProtocol))
}
