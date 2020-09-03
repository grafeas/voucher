package grafeasos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// isAttestationExistsErr returns true if the passed Error is an "AlreadyExists" gRPC error.
func isAttestionExistsErr(err error) bool {
	if nil == err {
		return false
	}

	return (codes.AlreadyExists == status.Code(err))
}
