package voucher

import (
	"fmt"

	"github.com/docker/distribution/reference"
)

// NewImageReference creates a new reference item with the passed URL as
// a reference to the target image.
func NewImageReference(url string) (reference.Canonical, error) {
	var imageData reference.Canonical

	rawRef, err := reference.Parse(url)
	if nil != err {
		return imageData, fmt.Errorf("can't use URL \"%s\" as image reference: %s", url, err)
	}

	canonicalRef, isCanonical := rawRef.(reference.Canonical)
	if !isCanonical {
		return imageData, fmt.Errorf("reference %s has no digest", rawRef.String())
	}

	return canonicalRef, nil
}
