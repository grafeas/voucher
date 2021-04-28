package voucher

import (
	"context"

	"github.com/docker/distribution/reference"
)

// Interface represents an interface to the Voucher API. Typically Voucher API
// clients would implement it.
type Interface interface {
	Check(ctx context.Context, check string, image reference.Canonical) (Response, error)
	Verify(ctx context.Context, check string, image reference.Canonical) (Response, error)
}
