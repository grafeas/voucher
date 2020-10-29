package voucher

import (
	"context"
	"fmt"

	"github.com/grafeas/voucher/repository"
	"github.com/docker/distribution/reference"
)

// MetadataClient is an interface that represents something that communicates
// with the Metadata server.
type MetadataClient interface {
	CanAttest() bool
	NewPayloadBody(ImageData) (string, error)
	GetVulnerabilities(context.Context, ImageData) ([]Vulnerability, error)
	GetBuildDetail(context.Context, reference.Canonical) (repository.BuildDetail, error)
	AddAttestationToImage(context.Context, ImageData, Attestation) (SignedAttestation, error)
	GetAttestations(context.Context, ImageData) ([]SignedAttestation, error)
	Close()
}

// NoMetadataError is an error that is returned when we request metadata that
// should exist but doesn't. It's a general error that will wrap more specific
// errors if desired.
type NoMetadataError struct {
	Type MetadataType
	Err  error
}

// Error returns the error value of this NoMetadataError as a string.
func (err *NoMetadataError) Error() string {
	return fmt.Sprintf("no metadata of type %s returned: %s", err.Type, err.Err)
}

// IsNoMetadataError returns true if the passed error is a NoMetadataError.
func IsNoMetadataError(err error) bool {
	_, ok := err.(*NoMetadataError)
	return ok
}
