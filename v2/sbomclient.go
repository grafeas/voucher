package voucher

import (
	"context"
	"errors"

	"github.com/docker/distribution/reference"
)

type SBOM struct {
}

// SBOMClient is an interface that represents something that gets SBOMs
type SBOMClient interface {
	GetVulnerabilities(ctx context.Context, ref reference.Canonical) ([]Vulnerability, error)
	GetSBOM(context.Context, reference.Canonical) (SBOM, error)
	Close()
}

// NoSBOMError is an error that is returned when we request sbom that should exist but doesn't
var NoSBOMError = errors.New("No SBOM was found")
