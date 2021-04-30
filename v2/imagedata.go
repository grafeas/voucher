package voucher

import (
	"fmt"

	"github.com/docker/distribution/reference"
)

// ImageData is a Canonical Reference to the Image (includes digest and URL).
type ImageData = reference.Canonical

// NewImageData creates a new ImageData item with the passed URL as
// a reference to the target image.
func NewImageData(url string) (ImageData, error) {
	var imageData ImageData
	rawRef, err := reference.Parse(url)
	if nil != err {
		return imageData, fmt.Errorf("can't use URL in ImageData: %s", err)
	}

	canonicalRef, isCanonical := rawRef.(reference.Canonical)
	if !isCanonical {
		return imageData, fmt.Errorf("reference %s has no digest", rawRef.String())
	}

	return canonicalRef, nil
}
