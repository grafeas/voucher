package docker

import (
	"errors"
	"net/http"

	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"

	"github.com/grafeas/voucher/docker/uri"
)

// GetDigestFromTagged gets an image's digest from the passed tag.
// Returns a digest.Digest, or an error.
func GetDigestFromTagged(client *http.Client, image reference.NamedTagged) (digest.Digest, error) {
	blank := digest.Digest("")

	request, err := http.NewRequest(http.MethodHead, uri.GetTagManifestURI(image), nil)
	if err != nil {
		return blank, err
	}

	request.Header.Add("Accept", schema2.MediaTypeManifest)
	request.Header.Add("Accept", schema1.MediaTypeManifest)
	request.Header.Add("Accept", schema1.MediaTypeSignedManifest)

	resp, err := client.Do(request)
	if err != nil {
		return blank, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return blank, responseToError(resp)
	}

	imageDigest := digest.Digest(resp.Header.Get("Docker-Content-Digest"))
	if "" == string(imageDigest) {
		return blank, errors.New("empty digest returned for image")
	}

	return imageDigest, nil
}

// GetDigestFromCanonical gets an image's digest from the passed canonical
// image. Returns a digest.Digest, or an error.
func GetDigestFromCanonical(client *http.Client, image reference.Canonical) (digest.Digest, error) {
	blank := digest.Digest("")

	request, err := http.NewRequest(http.MethodHead, uri.GetDigestManifestURI(image), nil)
	if err != nil {
		return blank, err
	}

	request.Header.Add("Accept", schema2.MediaTypeManifest)
	request.Header.Add("Accept", schema1.MediaTypeManifest)
	request.Header.Add("Accept", schema1.MediaTypeSignedManifest)

	resp, err := client.Do(request)
	if err != nil {
		return blank, err
	}

	if resp.StatusCode >= 300 {
		return blank, responseToError(resp)
	}

	_ = resp.Body.Close()

	imageDigest := digest.Digest(resp.Header.Get("Docker-Content-Digest"))
	if "" == string(imageDigest) {
		return blank, errors.New("empty digest returned for image")
	}

	return imageDigest, nil
}
