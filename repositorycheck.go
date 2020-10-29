package voucher

import "github.com/grafeas/voucher/repository"

// RepositoryCheck represents a Voucher check that needs to lookup
// information about an image from the repository that it's source code
// is stored in.
//
// RepositoryCheck implements a MetadataCheck, as containers normally
// do not contain information about their source repositories. This
// enables us to take advantage of Grafeas (or other metadata systems)
// which track build information for an image, in addition to signatures
// and (possibly) vulnerability information.
type RepositoryCheck interface {
	MetadataCheck
	SetRepositoryClient(repository.Client)
}
