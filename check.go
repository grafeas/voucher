package voucher

import (
	"context"
	"errors"

	"github.com/docker/distribution/reference"
)

// ErrNoCheck is an error that is returned when a requested check hasn't
// been registered.
var ErrNoCheck = errors.New("requested check doesn't exist")

// Check represents a Voucher test.
type Check interface {
	Check(context.Context, reference.Canonical) (bool, error)
}
