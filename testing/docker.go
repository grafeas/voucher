package vtesting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/libtrust"
)

// RateLimitOutput is the data that is returned when we similuate a Docker
// registry call that has been rate limited.
const RateLimitOutput = "<html><body>Rate Limited</body></html>"

// dockerAPIMock mocks the Docker API.
type dockerAPIMock struct {
	privateKey libtrust.PrivateKey
}

// ServeHTTP implements the http.Handler interface, responding to valid requests with good data, and
// invalid requests with garbage data.
func (mock *dockerAPIMock) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/v2/path/to/image/manifests/latest", "/v2/path/to/image/manifests/sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da":
		writer.Header().Set("Docker-Content-Digest", "sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da")
		mimeType, raw, _ := NewTestManifest().Payload()
		rawRespond(writer, mimeType, string(raw))
		return
	case "/v2/schema1image/manifests/latest", "/v2/schema1image/manifests/sha256:03f65aeeb2e8e8db022b297cae4cdce9248633f551452e63ba520d1f9ef2eca0":
		writer.Header().Set("Docker-Content-Digest", "sha256:03f65aeeb2e8e8db022b297cae4cdce9248633f551452e63ba520d1f9ef2eca0")
		jsonRespond(writer, schema1.MediaTypeManifest, NewTestSchema1Manifest())
		return
	case "/v2/schema1imagesigned/manifests/latest", "/v2/schema1imagesigned/manifests/sha256:18e6e7971438ab792d13563dcd8972acf4445bc0dcfdff84a6374d63a9c3ed62":
		writer.Header().Set("Docker-Content-Digest", "sha256:18e6e7971438ab792d13563dcd8972acf4445bc0dcfdff84a6374d63a9c3ed62")
		mimeType, raw, _ := NewTestSchema1SignedManifest(mock.privateKey).Payload()
		rawRespond(writer, mimeType, string(raw))
		return
	case "/v2/path/to/ratelimited/manifests/latest", "/v2/path/to/ratelimited/manifests/sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da":
		rawRespond(writer, "text/html", RateLimitOutput)
		return
	case "/v2/path/to/image/blobs/sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7":
		jsonRespond(writer, schema2.MediaTypeImageConfig, NewTestImageConfig())
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

	handler.privateKey = NewPrivateKey()

	server := httptest.NewTLSServer(handler)
	return server
}

// jsonRespond wraps the appropriate http.ResponseWriter calls to return
// JSON encoded data to the testing client. If an error occurs, it will call
// http.Error on the writer.
func jsonRespond(writer http.ResponseWriter, content string, v interface{}) {
	writer.Header().Set("Content-Type", content)
	err := json.NewEncoder(writer).Encode(v)
	if nil != err {
		http.Error(writer, fmt.Sprintf("failed to handle request: %s", err), 500)
	}
}

// rawRespond wraps the appropriate http.ResponseWriter calls to return raw
// (unmarshaled) data to the testing client.  If an error occurs, it will call
// http.Error on the writer.
func rawRespond(writer http.ResponseWriter, content, body string) {
	writer.Header().Set("Content-Type", content)
	_, err := fmt.Fprintln(writer, body)
	if nil != err {
		http.Error(writer, fmt.Sprintf("failed to handle request: %s", err), 500)
	}
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
