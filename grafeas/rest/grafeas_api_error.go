package rest

import "fmt"

//GrafeasAPIError to store grafeas API errors
type GrafeasAPIError struct {
	statusCode  int
	url         string
	method      string
	requestData string
}

// Error returns the grafeas API error as a string.
func (err *GrafeasAPIError) Error() string {
	if err.requestData != "" {
		return fmt.Sprintf("error getting REST data with status code %d for url %s and method %s with data: %v", err.statusCode, err.url, err.method, err.requestData)
	}
	return fmt.Sprintf("error getting REST data with status code %d for url %s and method %s", err.statusCode, err.url, err.method)
}

// NewGrafeasAPIError creates a new GrafeasAPIError
func NewGrafeasAPIError(statusCode int, url, method string, data []byte) error {
	return &GrafeasAPIError{
		statusCode:  statusCode,
		url:         url,
		method:      method,
		requestData: string(data),
	}
}
