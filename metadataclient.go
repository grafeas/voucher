package voucher

import (
	"context"
	"github.com/docker/distribution/reference"
)

// MetadataClient is an interface that represents something that communicates
// with the Metadata server.
type MetadataClient interface {
	CanAttest() bool
	NewPayloadBody(ImageData) (string, error)
	GetMetadata(context.Context, ImageData, MetadataType) ([]MetadataItem, error)
	GetBuildDetails(context.Context, reference.Canonical) ([]BuildDetail, error)
	AddAttestationToImage(context.Context, ImageData, AttestationPayload) (MetadataItem, error)
	Close()
}
