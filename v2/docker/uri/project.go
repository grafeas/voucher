package uri

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
)

type ErrNoProjectInReference struct {
	ref reference.Reference
}

func (err *ErrNoProjectInReference) Error() string {
	return fmt.Sprintf("could not find project path in reference \"%s\"", err.ref)
}

// ReferenceToProjectName returns what should be the GCR project name for an
// image reference.
//
// For example, if an image is in the project "my-cool-project" the image path
// should start with `gcr.io/my-cool-project`.
func ReferenceToProjectName(ref reference.Reference) (string, error) {
	values := strings.Split(ref.String(), "/")
	if 2 < len(values) {
		if values[0] == "gcr.io" {
			return values[1], nil
		}
	}

	return "", &ErrNoProjectInReference{
		ref: ref,
	}
}
