package docker

import (
	"errors"
	"io"
	"net/http"
)

// responseToError converts the body of a response to an error.
func responseToError(resp *http.Response) error {
	b, err := io.ReadAll(resp.Body)
	if nil == err {
		err = errors.New("failed to load resource with status \"" + resp.Status + "\": " + string(b))
	}

	return errors.New("failed to load resource with error: " + err.Error())
}
