package sbomgcr

import (
	"context"

	"github.com/docker/docker/reference"
	sbom "github.com/grafeas/voucher/v2"
	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/signer"
)

// Client connects to GCR
type Client struct {
}

// GetVulnerabilities returns the detected vulnerabilities for the Image described by voucher.ImageData.
func (c *Client) GetVulnerabilities(ctx context.Context, ref reference.Canonical) (vulnerabilities []voucher.Vulnerability, err error) {
	return []voucher.Vulnerability{}, nil
}

// GetBuildDetail gets the BuildDetail for the passed image.
func (g *Client) GetSBOM(ctx context.Context, ref reference.Canonical) (sbom.SBOM, error) {
	return sbom.SBOM{}, nil
}

// NewClient creates a new containeranalysis Grafeas Client.
func NewClient(ctx context.Context, binauthProject string, keyring signer.AttestationSigner) (*Client, error) {
	client := &Client{}

	return client, nil
}
