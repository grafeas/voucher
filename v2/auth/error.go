package auth

import (
	"fmt"

	"github.com/docker/distribution/reference"
)

// Error is the error returned when authenticating while pulling manifests connecting to protected systems
type Error struct {
	Reason    string
	ImageName reference.Named
}

// Error returns a string for logging purposes
func (e *Error) Error() string {
	return fmt.Sprintf("auth failed: %s for %s", e.Reason, e.ImageName)
}

// NewAuthError returns an AuthError struct with the reason for erroring, as well as the Name of the image reference
func NewAuthError(reason string, imageName reference.Named) error {
	return &Error{reason, imageName}
}
