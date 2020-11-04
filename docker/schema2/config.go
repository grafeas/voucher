package schema2

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	dockerTypes "github.com/docker/docker/api/types"

	"github.com/grafeas/voucher/docker/uri"
)

type v2Blob struct {
	Config dockerTypes.ExecConfig `json:"container_config"`
}

// RequestConfig requests an image configuration from the server, based on the passed digest.
// Returns an ImageConfig or an error.
func RequestConfig(client *http.Client, ref reference.Canonical, manifest distribution.Manifest) (*dockerTypes.ExecConfig, error) {
	if !IsManifest(manifest) {
		return nil, errors.New("cannot request schema2 config for non-schema2 manifest")
	}

	v2Manifest := ToManifest(manifest)

	var wrapper v2Blob

	request, err := http.NewRequest(
		http.MethodGet,
		uri.GetBlobURI(ref, v2Manifest.Config.Digest),
		nil,
	)
	if nil != err {
		return nil, err
	}

	request.Header.Add("Accept", v2Manifest.Config.MediaType)

	resp, err := client.Do(request)
	if nil != err {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 300 {
		err = json.NewDecoder(resp.Body).Decode(&wrapper)
		if nil == err {
			return &wrapper.Config, nil
		}
	}

	return nil, err
}
