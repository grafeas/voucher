package gcr

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
)

// ReferenceToProjectName returns what should be the GCR project name for an
// image reference.
//
// For example, if an image is in the project "my-cool-project" the image path
// should start with `gcr.io/my-cool-project`.
func ReferenceToProjectName(ref reference.Reference) (string, error) {
	values := strings.Split(ref.String(), "/")
	if len(values) > 2 {
		if values[0] == "gcr.io" {
			return values[1], nil
		}
		if strings.HasSuffix(values[0], ".pkg.dev") {
			return values[1], nil
		}
	}

	return "", fmt.Errorf("could not find project path in reference %q", ref)
}
