package sbomgcr

import (
	"context"
	"fmt"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/docker/distribution/reference"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	gcr "github.com/google/go-containerregistry/pkg/v1/google"
	voucher "github.com/grafeas/voucher/v2"
)

// Client connects to GCR
type Client struct {
	authenticator authn.Authenticator
}

// GetVulnerabilities returns the detected vulnerabilities for the Image described by voucher.ImageData.
func (c *Client) GetVulnerabilities(ctx context.Context, ref reference.Canonical) (vulnerabilities []voucher.Vulnerability, err error) {
	return []voucher.Vulnerability{}, nil
}

// GetSBOM gets the SBOM for the passed image.
func (c *Client) GetSBOM(ctx context.Context, ref reference.Canonical) (cyclonedx.BOM, error) {
	digest := ref.Digest().String()
	tag := strings.Replace(digest, "@", "-", 1)

	// we can call the GetRepoTags method to get back all the gcr.Tags on a specific repo
	// then match the tag we make from the digest against it
	// then get the sbom

	return cyclonedx.BOM{}, nil
}

// GetRepoTags gets the gcr tags for the passed image.
func (c *Client) GetRepoTags(ctx context.Context, repoName string) (*gcr.Tags, error) {
	repository, err := name.NewRepository(repoName)
	if err != nil {
		fmt.Errorf("error returning authenticator %w", err)
	}

	tags, err := gcr.List(repository, gcr.WithAuth(c.authenticator), gcr.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error occured when trying to retrieve tags from this repo %w", err)
	}

	return tags, err
}

// NewClient creates a new
func NewClient() (*Client, error) {
	// This authenticator is the only one that seems to use ADC?? lets try it
	auth, err := gcr.NewEnvAuthenticator()
	if err != nil {
		fmt.Errorf("error returning authenticator %w", err)
	}

	client := &Client{authenticator: auth}

	return client, nil
}
