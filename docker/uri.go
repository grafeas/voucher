package docker

import (
	"fmt"
	"net/url"

	"github.com/docker/distribution/reference"
	digest "github.com/opencontainers/go-digest"
)

// GetTokenURI gets the token URI for the passed repository.
func GetTokenURI(ref reference.Named) string {
	hostname := reference.Domain(ref)
	repository := reference.Path(ref)

	query := url.Values{}
	query.Set("service", hostname)
	query.Set("scope", "repository:"+repository+":*")

	return fmt.Sprintf("https://%s/v2/token?%s", hostname, query.Encode())
}

// GetBlobURI gets a blob URI based on the passed repository and
// digest.
func GetBlobURI(ref reference.Named, digest digest.Digest) string {
	hostname := reference.Domain(ref)
	repository := reference.Path(ref)

	return fmt.Sprintf("https://%s/v2/%s/blobs/%s", hostname, repository, digest)
}

// GetManifestURI gets a manifest URI based on the passed repository and
// digest.
func GetManifestURI(ref reference.Canonical) string {
	hostname := reference.Domain(ref)
	repository := reference.Path(ref)

	return fmt.Sprintf("https://%s/v2/%s/manifests/%s", hostname, repository, ref.Digest())
}

// GetTagManifestURI gets a manifest URI based on the passed repository and
// tag.
func GetTagManifestURI(ref reference.NamedTagged) string {
	hostname := reference.Domain(ref)
	repository := reference.Path(ref)

	return fmt.Sprintf("https://%s/v2/%s/manifests/%s", hostname, repository, ref.Tag())
}
