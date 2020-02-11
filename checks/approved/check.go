package approved

import (
	"context"
	"errors"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/repository"
)

// ErrNoBuildData is an error returned if we can't pull any BuildData from
// Grafeas for an image.
var ErrNoBuildData = errors.New("no build metadata associated with this image")

// ErrNeedsRepositoryClient is an error returned if there is no repository client
// configured for an image
var ErrNeedsRepositoryClient = errors.New("this check requires a repository client")

type check struct {
	metadataClient   voucher.MetadataClient
	repositoryClient repository.Client
}

// SetMetadataClient sets the MetadataClient for this Check.
func (g *check) SetMetadataClient(metadataClient voucher.MetadataClient) {
	g.metadataClient = metadataClient
}

// SetRepositoryClient sets the repository.Client for this Check.
func (g *check) SetRepositoryClient(repositoryClient repository.Client) {
	g.repositoryClient = repositoryClient
}

// Check checks that the code used to built the image passed all required checks from its source repository
func (g *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	buildDetail, err := g.metadataClient.GetBuildDetail(ctx, i)
	if err != nil {
		if voucher.IsNoMetadataError(err) {
			return false, ErrNoBuildData
		}
		return false, err
	}

	if g.repositoryClient == nil {
		return false, ErrNeedsRepositoryClient
	}

	commit, err := g.repositoryClient.GetCommit(ctx, buildDetail)
	if nil != err {
		return false, err
	}

	defaultBranch, err := g.repositoryClient.GetDefaultBranch(ctx, buildDetail)
	if nil != err {
		return false, err
	}

	if !isFromBranch(defaultBranch, commit) ||
		!isSigned(commit) ||
		!isMergeCommit(commit, defaultBranch) ||
		!passedCI(commit) {
		return false, nil
	}

	return true, nil
}

// isFromBranch checks that the commit is the most recent commit on the branch
func isFromBranch(branch repository.Branch, commit repository.Commit) bool {
	return commit.URL == branch.CommitRefs[0].URL
}

// isSigned checks that the commit is signed
func isSigned(commit repository.Commit) bool {
	return commit.IsSigned
}

// isMergeCommit checks that the commit is a merge commit
func isMergeCommit(commit repository.Commit, branch repository.Branch) bool {
	for _, pullRequest := range commit.AssociatedPullRequests {
		if pullRequest.IsMerged && pullRequest.MergeCommit.URL == commit.URL {
			return true
		}
	}
	return false
}

// passedCI checks that all Github CI checks have completed and passed
func passedCI(commit repository.Commit) bool {
	return commit.Status == repository.CommitStatusSuccess
}

func init() {
	voucher.RegisterCheckFactory("approved", func() voucher.Check {
		return new(check)
	})
}
