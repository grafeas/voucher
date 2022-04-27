package sbom

import (
	"context"

	"github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/sbomgcr"
)

// check is a check that verifies if there's an sbom attached with
// the container image
type check struct {
	sbomClient voucher.SBOMClient
}

// SetSBOMClient sets the sbom / gcr client that this check will use
// for its run.
func (c *check) SetSBOMClient(sbomClient voucher.SBOMClient) {
	c.sbomClient = sbomClient
}

// hasSBOM returns true if the passed image has an SBOM attached
func (c *check) hasSBOM(i voucher.ImageData) bool {
	// Parse the image reference
	imageName := i.Name()
	tag := sbomgcr.GetSBOMTagFromImage(i)

	// Get digest of the sbom and build a reference string
	// So we can pull the sbom from the image repository
	sbomDigest, err := c.sbomClient.GetSBOMDigestWithTag(context.Background(), imageName, tag)
	if err != nil {
		return false
	}

	sbomName := imageName + "@" + sbomDigest
	_, err = c.sbomClient.GetSBOM(context.Background(), sbomName)
	return err == nil
}

// check checks if an image was built by a trusted source
func (c *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	if !c.hasSBOM(i) {
		return false, voucher.ErrNoSBOM
	}
	// add more
	return true, nil
}

func init() {
	voucher.RegisterCheckFactory("sbom", func() voucher.Check {
		return new(check)
	})
}
