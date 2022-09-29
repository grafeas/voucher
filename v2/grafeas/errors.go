package grafeas

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Grafeas client errors
var (
	errNoOccurrences         = errors.New("no occurrences returned for image")
	errDiscoveriesUnfinished = errors.New("discoveries have not finished processing")
)

// isAttestationExistsErr returns true if the passed Error is an "AlreadyExists" gRPC error.
func isAttestationExistsErr(err error) bool {
	if nil == err {
		return false
	}

	return (codes.AlreadyExists == status.Code(err))
}
