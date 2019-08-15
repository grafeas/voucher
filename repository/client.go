package repository

import "context"

// Client is a client for a version control source
type Client interface {
	GetCommitInfo(ctx context.Context, commitURL string) (CommitInfo, error)
	GetOrganization(ctx context.Context, uri string) (Organization, error)
	GetDefaultBranch(ctx context.Context, commitURL string) (DefaultBranch, error)
}
