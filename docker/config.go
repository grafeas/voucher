package docker

import (
	"errors"
	"net/http"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	digest "github.com/opencontainers/go-digest"
)

// RequestImageConfig requests an image configuration from the server, based on the passed
// reference. Returns an ImageConfig or an error.
func RequestImageConfig(client *http.Client, ref reference.Canonical) (ImageConfig, error) {

	manifest, err := RequestManifest(client, ref)
	if nil != err {
		return ImageConfig{}, err
	}

	if "" == manifest.Config.Digest {
		return ImageConfig{}, errors.New("image does not have any configuration")
	}

	return RequestConfig(client, ref, manifest.Config.Digest)
}

// RequestConfig requests an image configuration from the server, based on the passed digest.
// Returns an ImageConfig or an error.
func RequestConfig(client *http.Client, ref reference.Named, digest digest.Digest) (ImageConfig, error) {
	var config ImageConfig

	request, err := http.NewRequest(http.MethodGet, GetBlobURI(ref, digest), nil)
	if nil != err {
		return config, err
	}

	request.Header.Add("Accept", schema2.MediaTypeImageConfig)

	err = doDockerCall(client, request, &config)

	return config, err
}
