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

	"github.com/grafeas/voucher/cmd/config"
	"github.com/grafeas/voucher/metrics"
)

var testParams = []byte(`
{  
	"image_url":"gcr.io/somewhere/image@sha256:cb749360c5198a55859a7f335de3cf4e2f64b60886a2098684a2f9c7ffca81f2"
}
`)

var server *Server

const testUsername = "vouchertester"
const testPassword = "testingvoucher"
const testPasswordHash = "$2a$10$.PaOjV8GdqSHSmUtfolsJeF6LsAq/3CNsFCYGb3IoN/mO9xj1c/yG"

func TestMain(m *testing.M) {
	config.FileName = "../testdata/config.toml"

	config.InitConfig()

	serverConfig := &Config{
		Port:        viper.GetInt("server.port"),
		Timeout:     viper.GetInt("server.timeout"),
		RequireAuth: true,
		Username:    testUsername,
		PassHash:    testPasswordHash,
	}
	secrets, _ := config.ReadSecrets()
	server = NewServer(serverConfig, secrets, &metrics.NoopClient{})

	for groupName, checks := range config.GetRequiredChecksFromConfig() {
		server.SetCheckGroup(groupName, checks)
	}

	os.Exit(m.Run())
}

func TestGoodAuthentication(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/all", bytes.NewReader(testParams))
	require.NoError(t, err)

	router := NewRouter(server)

	req.SetBasicAuth(testUsername, testPassword)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestBadAuthentication(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/all", bytes.NewReader(testParams))
	require.NoError(t, err)

	router := NewRouter(server)

	req.SetBasicAuth(testUsername, "not the password")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check if the status code is not 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`
		{  
			image_url:poorly-formatted-json,
		}
		`)

	router := NewRouter(server)

	// Use the check groups configured in the test config
	config.FileName = "../testdata/config.toml"
	config.InitConfig()

	for _, route := range getRoutes(server) {
		path := route.Path

		if individualCheckPath == path {
			path = "/diy"
		} else if verifyCheckPath == path {
			path = "/diy/verify"
		} else if healthCheckPath == path {
			continue
		}

		req, err := http.NewRequest(route.Method, path, bytes.NewReader(invalidJSON))
		require.NoError(t, err)

		req.SetBasicAuth(testUsername, testPassword)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equalf(t, http.StatusUnprocessableEntity, recorder.Code, "failed to get Unprocessable Entity status on %s", path)
	}
}

func TestHandlerStatus(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, healthCheckPath, nil)
	require.NoError(t, err)

	router := NewRouter(server)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code, "handler for health check failed")
}
