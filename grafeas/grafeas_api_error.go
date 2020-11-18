package grafeas

import "fmt"

//APIError to store grafeas API errors
type APIError struct {
	statusCode  int
	url         string
	method      string
	requestData string
}

// Error returns the grafeas API error as a string.
func (err *APIError) Error() string {
	if err.requestData != "" {
		return fmt.Sprintf("error getting REST data with status code %d for url %s and method %s with data: %v", err.statusCode, err.url, err.method, err.requestData)
	}
	return fmt.Sprintf("error getting REST data with status code %d for url %s and method %s", err.statusCode, err.url, err.method)
}

// NewAPIError creates a new APIError
func NewAPIError(statusCode int, url, method string, data []byte) error {
	return &APIError{
		statusCode:  statusCode,
		url:         url,
		method:      method,
		requestData: string(data),
	}
}
