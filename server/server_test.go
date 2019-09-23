package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Shopify/voucher/cmd/config"
)

var testParams = []byte(`
{  
	"image_url":"gcr.io/somewhere/image@sha256:cb749360c5198a55859a7f335de3cf4e2f64b60886a2098684a2f9c7ffca81f2"
}
`)

const testUsername = "vouchertester"
const testPassword = "testingvoucher"
const testPasswordHash = "$2a$10$.PaOjV8GdqSHSmUtfolsJeF6LsAq/3CNsFCYGb3IoN/mO9xj1c/yG"

func TestMain(m *testing.M) {
	config.FileName = "../testdata/config.toml"

	config.InitConfig()

	serverConfig = &Config{
		Port:        viper.GetInt("server.port"),
		Timeout:     viper.GetInt("server.timeout"),
		RequireAuth: true,
		Username:    testUsername,
		PassHash:    testPasswordHash,
	}

	os.Exit(m.Run())
}

func TestGoodAuthentication(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/all", bytes.NewReader(testParams))
	require.NoError(t, err)

	req.SetBasicAuth(testUsername, testPassword)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleAll)

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestBadAuthentication(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/all", bytes.NewReader(testParams))
	require.NoError(t, err)

	router := NewRouter()

	req.SetBasicAuth(testUsername, "not the password")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check if the status code is not 401 Unauthorized
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`
		{  
			image_url:poorly-formatted-json,
		}
		`)

	router := NewRouter()

	for _, route := range Routes {
		path := route.Path
		if "/{check}" == path {
			path = "/diy"
		} else if healthCheckPath == path {
			continue
		}

		req, err := http.NewRequest(route.Method, path, bytes.NewReader(invalidJSON))
		require.NoError(t, err)

		req.SetBasicAuth(testUsername, testPassword)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, recorder.Code, http.StatusUnprocessableEntity)
	}
}

func TestHandlerStatus(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, healthCheckPath, nil)
	require.NoError(t, err)

	router := NewRouter()

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, recorder.Code, http.StatusOK, "handler for health check failed")
}
