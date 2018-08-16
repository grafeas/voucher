package voucher

import (
	"errors"
)

// ErrNoCheck is an error that is returned when a requested check hasn't
// been registered.
var ErrNoCheck = errors.New("requested check doesn't exist")

// Check represents a Voucher test.
type Check interface {
	Check(ImageData) (bool, error)
}
