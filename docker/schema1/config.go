package schema1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	dockerTypes "github.com/docker/docker/api/types"
)

type v1Blob struct {
	Config dockerTypes.ExecConfig `json:"config"`
}

// RequestConfig retrieves the manifest from the associated configuration.
// Unlike in v2 manifests, v1 manifests have the configuration stored in the
// history, so we can safely ignore the http.Client passed to this function.
func RequestConfig(_ *http.Client, _ reference.Canonical, manifest distribution.Manifest) (*dockerTypes.ExecConfig, error) {
	if !IsManifest(manifest) {
		return nil, errors.New("cannot request schema1 config for non-schema1 manifest")
	}

	v1Manifest := ToManifest(manifest)

	if len(v1Manifest.History) < 1 {
		return nil, errors.New("no history in manifest")
	}

	configBlob := v1Blob{}

	err := json.Unmarshal([]byte(v1Manifest.History[0].V1Compatibility), &configBlob)
	if nil != err {
		return nil, err
	}

	return &configBlob.Config, nil
}
