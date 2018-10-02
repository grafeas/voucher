package vtesting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	dockerTypes "github.com/docker/docker/api/types"
)

// dockerAPIMock mocks the Docker API.
type dockerAPIMock struct {
}

// ServeHTTP implements the http.Handler interface, responding to valid requests with good data, and
// invalid requests with garbage data.
func (mock *dockerAPIMock) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/v2/path/to/image/manifests/latest", "/v2/path/to/image/manifests/sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da":
		writer.Header().Set("Docker-Content-Digest", "sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da")
		respond(writer, schema2.MediaTypeManifest, NewTestManifest())
		return
	case "/v2/path/to/image/blobs/sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7":
		respond(writer, schema2.MediaTypeImageConfig, NewTestImageConfig())
		return
	case "/v2/path/to/bad/image/manifests/latest", "/v2/path/to/bad/image/manifests/sha256:bad8c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da":
		http.Error(writer, "image doesn't exist", 404)
		return
	}

	http.Error(writer, fmt.Sprintf("failed to handle request: %s", req.URL.Path), 500)
}

// NewTestDockerServer creates a new mock of the Docker registry
func NewTestDockerServer(t *testing.T) *httptest.Server {
	handler := new(dockerAPIMock)

	server := httptest.NewTLSServer(handler)
	return server
}

// respond wraps the appropriate http.ResponseWriter calls to return data to the testing client.
// If an error occurs, it will call http.Error on the writer.
func respond(writer http.ResponseWriter, content string, v interface{}) {
	writer.Header().Set("Content-Type", content)
	err := json.NewEncoder(writer).Encode(v)
	if nil != err {
		http.Error(writer, fmt.Sprintf("failed to handle request: %s", err), 500)
	}
}

// NewTestManifest creates a test manifest for our mock Docker API.
func NewTestManifest() schema2.Manifest {
	manifest := schema2.Manifest{
		Config: distribution.Descriptor{
			MediaType: schema2.MediaTypeImageConfig,
			Size:      7023,
			Digest:    "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7",
		},
		Layers: []distribution.Descriptor{
			{
				MediaType: schema2.MediaTypeLayer,
				Size:      32654,
				Digest:    "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
			},
			{
				MediaType: schema2.MediaTypeLayer,
				Size:      16724,
				Digest:    "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b",
			},
			{
				MediaType: schema2.MediaTypeLayer,
				Size:      73109,
				Digest:    "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736",
			},
		},
	}

	manifest.SchemaVersion = 2
	manifest.MediaType = schema2.MediaTypeManifest

	return manifest
}

// NewTestImageConfig creates a test image Config for our mock Docker API.
func NewTestImageConfig() interface{} {
	config := struct {
		ContainerConfig dockerTypes.ExecConfig `json:"container_config"`
	}{
		ContainerConfig: dockerTypes.ExecConfig{
			User: "root",
		},
	}

	return config
}
