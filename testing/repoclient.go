package vtesting

import (
	r "github.com/Shopify/voucher/repository"
)

// RepositoryClient is a client for a version control source
type RepositoryClient interface {
	r.Client
	AddRepository(org r.Organization, name string, commits map[string]r.Commit, branches map[string]r.Branch)
}
