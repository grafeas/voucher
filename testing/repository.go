package vtesting

import (
	"github.com/Shopify/voucher/repository"
)

// Repository contains repository information pertaining to a repository
type Repository struct {
	Org      repository.Organization
	Name     string
	Commits  map[string]repository.Commit
	Branches map[string]repository.Branch
}
