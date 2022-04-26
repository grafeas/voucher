package sbom

import (
	"context"
	"errors"

	"github.com/grafeas/voucher/v2"
)

// check is a check that verifies if there's an sbom attached with
// the container image
type check struct {
	sbomClient voucher.SBOMClient
}

// ErrNoSBOMFound  is returned when an image does not have
// any sboms attached to it
var ErrNoSBOMFound = errors.New("image has no sbom attached")

// SetSBOMClient sets the sbom / gcr client that this check will use
// for its run.
func (c *check) SetSBOMClient(sbomClient voucher.SBOMClient) {
	c.sbomClient = sbomClient
}

// hasSBOM returns true if the passed image has an SBOM attached
func (c *check) hasSBOM(i voucher.ImageData) bool {
	_, err := c.sbomClient.GetSBOM(context.Background(), i)

	return err == nil
}

// check checks if an image was built by a trusted source
func (c *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	if !c.hasSBOM(i) {
		return false, ErrNoSBOMFound
	}
	// add more
	return true, nil
}

func init() {
	voucher.RegisterCheckFactory("sbom", func() voucher.Check {
		return new(check)
	})
}
