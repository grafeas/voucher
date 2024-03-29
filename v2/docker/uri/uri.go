package uri

import (
	"bytes"
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

	u := createURL(ref, "token")
	u.RawQuery = query.Encode()

	return u.String()
}

// GetBlobURI gets a blob URI based on the passed repository and
// digest.
func GetBlobURI(ref reference.Named, digest digest.Digest) string {
	u := createURL(ref, reference.Path(ref), "blobs", string(digest))
	return u.String()
}

// GetManifestURI gets a manifest URI based on the passed repository and label (tag or digest).
func GetManifestURI(ref reference.Named, label string) string {
	u := createURL(ref, reference.Path(ref), "manifests", label)
	return u.String()
}

// GetTagManifestURI gets a manifest URI based on the passed repository and
// tag.
func GetTagManifestURI(ref reference.NamedTagged) string {
	return GetManifestURI(ref, ref.Tag())
}

// GetDigestManifestURI gets a manifest URI based on the passed repository and
// digest.
func GetDigestManifestURI(ref reference.Canonical) string {
	return GetManifestURI(ref, string(ref.Digest()))
}

func createURL(ref reference.Named, pathSegments ...string) url.URL {
	hostname := reference.Domain(ref)

	var path bytes.Buffer
	path.WriteString("/v2")

	for _, pathSegment := range pathSegments {
		path.WriteString("/")
		path.WriteString(pathSegment)
	}

	var u url.URL
	u.Scheme = "https"
	u.Host = hostname
	u.Path = path.String()

	return u
}
