package auth

import (
	"fmt"

	"github.com/docker/distribution/reference"
)

// AuthError is the error returned when authenticating while pulling manifests connecting to protected systems
type AuthError struct {
	Reason    string
	ImageName reference.Named
}

// Error returns a string for logging purposes
func (e *AuthError) Error() string {
	return fmt.Sprintf("auth failed: %s for %s", e.Reason, e.ImageName)
}

// NewAuthError returns an AuthError struct with the reason for erroring, as well as the Name of the image reference
func NewAuthError(reason string, imageName reference.Named) error {
	return &AuthError{reason, imageName}
}
