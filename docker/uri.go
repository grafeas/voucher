package docker

import (
	"fmt"
	"net/url"

	"github.com/docker/distribution/reference"
)

// GetTokenURI gets the token URI for the passed repository.
func GetTokenURI(ref reference.Named) string {
	hostname, repository := reference.SplitHostname(ref)

	query := url.Values{}
	query.Set("service", hostname)
	query.Set("scope", "repository:"+repository+":*")

	return fmt.Sprintf("https://%s/v2/token?%s", hostname, query.Encode())
}

// GetBlobURI gets a blob URI based on the passed repository and
// digest.
func GetBlobURI(ref reference.Canonical) string {
	hostname, repository := reference.SplitHostname(ref)

	return fmt.Sprintf("https://%s/v2/%s/blobs/%s", hostname, repository, ref.Digest())
}

// GetManifestURI gets a manifest URI based on the passed repository and
// tag.
func GetManifestURI(ref reference.Canonical) string {
	hostname, repository := reference.SplitHostname(ref)

	return fmt.Sprintf("https://%s/v2/%s/manifests/%s", hostname, repository, ref.Digest())
}
