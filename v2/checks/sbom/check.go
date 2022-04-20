package sbom

import (
	"context"
	"errors"

	"github.com/grafeas/voucher/v2"
)

// There are two supported versions in
// https://github.com/spdx/tools-golang/tree/main/spdx
// i.e. 2.1 and 2.2; we'll need to support both for spdx

// check is a check that verifies if there's an sbom attached with
// the container image
type check struct {
	auth voucher.Auth
}

// ErrNoSBOMFound  is returned when an image does not have
// any sboms attached to it
var ErrNoSBOMFound = errors.New("image has no sbom attached")

// SetAuth sets the authentication system that this check will use
// for its run.
func (c *check) SetAuth(auth voucher.Auth) {
	c.auth = auth
}

// hasSBOM returns true if the passed image has an SBOM attached
func (c *check) hasSBOM(i voucher.ImageData) bool {
	// TODO: call GCR to check for the presence of an SBOM

	// digest := i.Digest()

	// TODO: query GCR to see if there's an OCI manifest whose tag
	// contains the digest

	return false
}

// check checks if an image was built by a trusted source
func (c *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	if !c.hasSBOM(i) {
		return false, ErrNoSBOMFound
	}

	return true, nil
}

func init() {
	voucher.RegisterCheckFactory("sbom", func() voucher.Check {
		return new(check)
	})
}
