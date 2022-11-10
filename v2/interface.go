package voucher

import (
	"context"
)

// Interface represents an interface to the Voucher API. Typically Voucher API
// clients would implement it.
type Interface interface {
	Check(ctx context.Context, check string, image string) (Response, error)
	Verify(ctx context.Context, check string, image string) (Response, error)
}
