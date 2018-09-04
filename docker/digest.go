package docker

import (
	"net/http"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
)

// GetDigestFromTagged gets an image's digest from the passed tag.
// Returns a digest.Digest, or an error.
func GetDigestFromTagged(client *http.Client, image reference.NamedTagged) (digest.Digest, error) {
	blank := digest.Digest("")

	request, err := http.NewRequest(http.MethodHead, GetTagManifestURI(image), nil)
	if err != nil {
		return blank, err
	}

	request.Header.Add("Accept", schema2.MediaTypeManifest)

	resp, err := client.Do(request)
	if err != nil {
		return blank, err
	}

	_ = resp.Body.Close()

	return digest.Digest(resp.Header.Get("Docker-Content-Digest")), nil

}
