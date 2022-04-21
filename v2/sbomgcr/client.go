package sbomgcr

import (
	"context"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/docker/distribution/reference"
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

// GetSBOM gets the SBOM for the passed image.
func (g *Client) GetSBOM(ctx context.Context, ref reference.Canonical) (cyclonedx.BOM, error) {
	digest := ref.Digest().String()
	tag := strings.Replace(digest, "@", "-", 1)
	// get repo (similar to gcloud describe) or parse as ReferenceToProjectName
	// lets authenticate this thing
	// list tags https://pkg.go.dev/github.com/google/go-containerregistry@v0.8.0/pkg/v1/google#List
	// match tag to sbom
	// get sbom

	return cyclonedx.BOM{}, nil
}

// NewClient creates a new
func NewClient(ctx context.Context, binauthProject string, keyring signer.AttestationSigner) (*Client, error) {
	client := &Client{}

	return client, nil
}
