package docker

import "fmt"

const (
	manifestType = "manifest"
	configType   = "config"
)

// APIError is a generic error structure representing a docker API call
// error. It tracks the type of request that failed, and either wraps an error
// or contains the body of an API call.
type APIError struct {
	callType      string
	requestStatus string
	requestBody   string
	err           error
}

// Error returns the docker API error as a string.
func (err *APIError) Error() string {
	if err.requestBody != "" {
		return fmt.Sprintf("failed to load %s with status %s: \"%s\"", err.callType, err.requestStatus, err.requestBody)
	}
	return fmt.Sprintf("failed to load %s: %s", err.callType, err.err)
}

// NewManifestError creates a new APIError specific to docker manifest requests.
// This version wraps the passed error.
func NewManifestError(err error) error {
	return &APIError{
		callType: manifestType,
		err:      err,
	}
}

// NewManifestErrorWithRequest creates a new APIError specific to docker
// manifest requests. This version wraps the passed HTTP response.
func NewManifestErrorWithRequest(status string, b []byte) error {
	return &APIError{
		callType:      manifestType,
		requestStatus: status,
		requestBody:   string(b),
	}
}

// NewConfigError creates a new APIError specific to docker config requests.
// This version wraps the passed error.
func NewConfigError(err error) error {
	return &APIError{
		callType: configType,
		err:      err,
	}
}

// NewConfigErrorWithRequest creates a new APIError specific to docker config
// requests. This version wraps the passed HTTP response.
func NewConfigErrorWithRequest(status string, b []byte) error {
	return &APIError{
		callType:      configType,
		requestStatus: status,
		requestBody:   string(b),
	}
}
