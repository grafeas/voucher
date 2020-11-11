package docker

import (
	"errors"
	"net/http"

	"github.com/docker/distribution/reference"
	dockerTypes "github.com/docker/docker/api/types"

	"github.com/grafeas/voucher/docker/schema1"
	"github.com/grafeas/voucher/docker/schema2"
)

// RequestImageConfig requests an image configuration from the server, based on the passed
// reference. Returns an ImageConfig or an error.
func RequestImageConfig(client *http.Client, ref reference.Canonical) (ImageConfig, error) {
	manifest, err := RequestManifest(client, ref)
	if nil != err {
		return nil, err
	}

	var config *dockerTypes.ExecConfig

	switch {
	case schema1.IsManifest(manifest):
		config, err = schema1.RequestConfig(client, ref, manifest)
	case schema2.IsManifest(manifest):
		config, err = schema2.RequestConfig(client, ref, manifest)
	default:
		err = errors.New("image does not have any configuration")
	}

	if nil != err {
		return nil, NewConfigError(err)
	}

	return &imageConfig{
		*config,
	}, nil
}
