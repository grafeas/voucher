package uri

import (
	"bytes"
	"net/url"

	"github.com/docker/distribution/reference"
	digest "github.com/opencontainers/go-digest"
)

// GetBlobURI gets a blob URI based on the passed repository and
// digest.
func GetBlobURI(ref reference.Named, digest digest.Digest) string {
	u := createURL(ref, reference.Path(ref), "blobs", string(digest))
	return u.String()
}

// GetManifestURI gets a manifest URI based on the passed repository and
// digest.
func GetManifestURI(ref reference.Canonical) string {
	u := createURL(ref, reference.Path(ref), "manifests", string(ref.Digest()))
	return u.String()
}

// GetTagManifestURI gets a manifest URI based on the passed repository and
// tag.
func GetTagManifestURI(ref reference.NamedTagged) string {
	u := createURL(ref, reference.Path(ref), "manifests", ref.Tag())
	return u.String()
}

// GetDigestManifestURI gets a manifest URI based on the passed repository and
// tag.
func GetDigestManifestURI(ref reference.Canonical) string {
	u := createURL(ref, reference.Path(ref), "manifests", string(ref.Digest()))
	return u.String()
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
