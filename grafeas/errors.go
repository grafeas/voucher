package grafeas

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Grafeas client errors
var (
	ErrNoOccurrences         = errors.New("no occurrences returned for image")
	ErrDiscoveriesUnfinished = errors.New("discoveries have not finished processing")
)

// IsAttestationExistsErr returns true if the passed Error is an "AlreadyExists" gRPC error.
func IsAttestationExistsErr(err error) bool {
	if nil == err {
		return false
	}

	return (codes.AlreadyExists == status.Code(err))
}
