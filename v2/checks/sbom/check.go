package sbom

import (
	"context"
	"strings"

	"github.com/grafeas/voucher/v2"
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
	tag := getSBOMTagFromImage(i)

	_, err := c.sbomClient.GetSBOM(context.Background(), imageName, tag)
	return err == nil
}

// GetSBOMTagFromImage returns the sbom tag from the image
func getSBOMTagFromImage(i voucher.ImageData) string {
	// Parse the image reference
	imageSHA := string(i.Digest())
	tag := strings.Replace(imageSHA, ":", "-", 1) + ".att"
	return tag
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
