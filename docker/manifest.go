package docker

import (
	"net/http"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher/docker/uri"
)

// RequestManifest requests an Manifest for the passed canonical image reference (an image URL
// with a digest specifying the built image). Returns a schema2.Manifest, or an error if
// there's an issue.
func RequestManifest(client *http.Client, ref reference.Canonical) (distribution.Manifest, error) {
	var manifest distribution.Manifest

	request, err := http.NewRequest(http.MethodGet, uri.GetManifestURI(ref), nil)
	if nil != err {
		return nil, err
	}

	request.Header.Add("Accept", schema2.MediaTypeManifest)
	request.Header.Add("Accept", schema1.MediaTypeManifest)
	request.Header.Add("Accept", schema1.MediaTypeSignedManifest)

	manifest, err = getDockerManifest(client, request)
	if nil != err {
		return nil, err
	}

	return manifest, nil
}
