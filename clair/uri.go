package clair

import (
	"path"
	"strings"

	digest "github.com/opencontainers/go-digest"
)

// GetNewLayerURI gets the new layer URI for the passed hostname.
func GetNewLayerURI(hostname string) string {
	return createURI(hostname, "v1/layers")
}

// GetLayerURI gets the layer URI for the passed digest on the passed hostname.
func GetLayerURI(hostname string, digest digest.Digest) string {
	return createURI(hostname, "v1/layers/", string(digest)) + "?vulnerabilities"
}

// createURI creates a new Clair URI based on the passed hostname. This will
// automatically add "https://" if a protocol is omitted.
func createURI(hostname string, pieces ...string) string {
	// if our URL doesn't have a protocol, we'll default to HTTPS.
	if !strings.HasPrefix(hostname, "https://") && !strings.HasPrefix(hostname, "http://") {
		hostname = "https://" + hostname
	}

	return hostname + "/" + path.Join(pieces...)
}
