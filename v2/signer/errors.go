package signer

import (
	"errors"
)

// ErrNoKeyForCheck is the error returned when Voucher does not have a key
// for the Check in question.
var ErrNoKeyForCheck = errors.New("no signing entity exists for check")
