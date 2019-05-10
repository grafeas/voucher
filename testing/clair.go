package vtesting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/coreos/clair/api/v1"
)

const (
	jsonContentType = "application/json;charset=utf-8"
	basicUsername   = "shopifolk"
	basicPassword   = "shopify"
)

type clairAPIMock struct {
	layerVulns map[string][]v1.Vulnerability
	vulns      map[string][]v1.Vulnerability
}

func (mock *clairAPIMock) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	user, pass, _ := req.BasicAuth()
	if user != basicUsername || pass != basicPassword {
		http.Error(writer, fmt.Sprintf("Unauthorized: %s", req.URL.Path), 401)
	}

	switch req.URL.Path {
	case "/v1/layers":
		var data map[string]v1.Layer
		_ = json.NewDecoder(req.Body).Decode(&data)

		current := data["Layer"].Name
		parent := data["Layer"].ParentName

		mock.layerVulns[current] = mock.layerVulns[parent]
		mock.layerVulns[current] = append(mock.layerVulns[current], mock.vulns[current]...)

		jsonRespond(writer, jsonContentType, "okay")
		return
	default:
		for key := range mock.layerVulns {
			if "/v1/layers/"+key == req.URL.Path {
				jsonRespond(writer, jsonContentType, createClairLayers(mock.layerVulns[key]))
				return
			}
		}
	}
	http.Error(writer, fmt.Sprintf("failed to handle request: %s", req.URL.Path), 500)
}

// NewTestClairServer creates a mock of Clair with a list of pre-defined clair
// vulnerabilities
func NewTestClairServer(t *testing.T, vulns map[string][]v1.Vulnerability) *httptest.Server {
	handler := new(clairAPIMock)
	handler.layerVulns = map[string][]v1.Vulnerability{}
	handler.vulns = vulns

	server := httptest.NewServer(handler)
	return server
}

func createClairLayers(clairVulns []v1.Vulnerability) v1.LayerEnvelope {
	layerEnv := v1.LayerEnvelope{
		Layer: &v1.Layer{
			Features: []v1.Feature{
				{
					Vulnerabilities: clairVulns,
				},
			},
		},
	}

	return layerEnv
}
