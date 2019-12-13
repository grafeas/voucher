package repository

import "context"

// Client is a client for a version control source
type Client interface {
	GetCommit(ctx context.Context, details BuildDetail) (Commit, error)
	GetOrganization(ctx context.Context, details BuildDetail) (Organization, error)
	GetBranch(ctx context.Context, details BuildDetail, name string) (Branch, error)
	GetDefaultBranch(ctx context.Context, details BuildDetail) (Branch, error)
}
