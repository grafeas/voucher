package voucher

import (
	"context"
	"errors"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/docker/distribution/reference"
)

// ErrNoSBOM is an error that is returned when we request sbom that should exist but doesn't
var ErrNoSBOM = errors.New("no SBOM was found")

// SBOMClient is an interface that represents something that gets SBOMs
type SBOMClient interface {
	GetVulnerabilities(ctx context.Context, ref reference.Canonical) ([]Vulnerability, error)
	GetSBOM(context.Context, string, string) (cyclonedx.BOM, error)
}
