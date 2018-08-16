package docker

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
)

// RequestImageConfig requests an image configuration from the server, based on the passed
// reference. Returns an ImageConfig or an error.
func RequestImageConfig(token OAuthToken, ref reference.Canonical) (ImageConfig, error) {

	manifest, err := RequestManifest(token, ref)
	if nil != err {
		return ImageConfig{}, err
	}

	if "" == manifest.Config.Digest {
		return ImageConfig{}, errors.New("image does not have any configuration")
	}

	configRef, err := reference.WithDigest(ref, manifest.Config.Digest)
	if nil != err {
		return ImageConfig{}, fmt.Errorf("could not create configuration reference: %s", err)
	}

	return RequestConfig(token, configRef)
}

// RequestConfig requests an image configuration from the server, based on the passed digest.
// Returns an ImageConfig or an error.
func RequestConfig(token OAuthToken, ref reference.Canonical) (ImageConfig, error) {
	var config ImageConfig

	request, err := http.NewRequest(http.MethodGet, GetBlobURI(ref), nil)
	if nil != err {
		return config, err
	}

	request.Header.Add("Accept", schema2.MediaTypeImageConfig)
	setBearerToken(request, token)

	err = doDockerCall(request, &config)

	return config, err
}
