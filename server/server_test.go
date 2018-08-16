package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Shopify/voucher/cmd/config"
)

var testParams = []byte(`
{  
	"image_url":"gcr.io/somewhere/image@sha256:cb749360c5198a55859a7f335de3cf4e2f64b60886a2098684a2f9c7ffca81f2",
	"project":"project"
}
`)

func TestMain(m *testing.M) {
	config.FileName = "../config/config.toml"

	config.InitConfig()

	os.Exit(m.Run())
}

func TestInvalidJSON(t *testing.T) {
	var invalidJSON = []byte(`
		{  
			image_url:poorly-formatted-json,
			project:"project
		}
		`)

	for _, route := range Routes {
		req, err := http.NewRequest(route.Method, route.Path, bytes.NewReader(invalidJSON))
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(route.HandlerFunc)

		handler.ServeHTTP(recorder, req)

		// Check the status code is 422 Unprocessable Entity
		if status := recorder.Code; status != http.StatusUnprocessableEntity {
			if healthCheckPath == route.Path {
				continue
			}
			t.Errorf("%s handler returned wrong status code: got %v wanted %v",
				route.Path, status, http.StatusUnprocessableEntity)
		}
	}
}

func TestHandlerStatus(t *testing.T) {

	for _, route := range Routes {
		req, err := http.NewRequest(route.Method, route.Path, bytes.NewReader(testParams))
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(route.HandlerFunc)

		handler.ServeHTTP(recorder, req)

		// Check the status code is what we expect
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler %v: returned wrong status code: got %v wanted %v",
				route.Name, status, http.StatusOK)
		}
	}
}
