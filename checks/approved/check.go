package approved

import (
	"context"
	"errors"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/repository"
)

// ErrNoBuildData is an error returned if we can't pull any BuildData from
// Grafeas for an image.
var ErrNoBuildData = errors.New("no build metadata associated with this image")

// ErrNeedsRepositoryClient is an error returned if there is no repository client
// configured for an image
var ErrNeedsRepositoryClient = errors.New("this check requires a repository client")

var ErrNotSigned = errors.New("commit was not signed by a valid key")
var ErrNotOnDefaultBranch = errors.New("commit is not the latest commit on the production branch")
var ErrNotMergeCommit = errors.New("commit is not a merge commit")
var ErrMissingRequiredApprovals = errors.New("the PR associated with this commit does not have the required number of approvals")
var ErrNotPassedCI = errors.New("commit did not pass CI in source code repository")

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

	if !isFromBranch(defaultBranch, commit) {
		return false, ErrNotOnDefaultBranch
	}

	if !isSigned(commit) {
		return false, ErrNotSigned
	}

	if result, reason := isApprovedMergeCommit(commit); !result {
		return result, reason
	}

	if !passedCI(commit) {
		return false, ErrNotPassedCI
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

// isApprovedMergeCommit checks that the commit is a merge commit
func isApprovedMergeCommit(commit repository.Commit) (passed bool, reason error) {
	for _, pullRequest := range commit.AssociatedPullRequests {
		if pullRequest.IsMerged && pullRequest.MergeCommit.URL == commit.URL {
			if !pullRequest.HasRequiredApprovals {
				return false, ErrMissingRequiredApprovals
			}
			return true, nil
		}
	}
	return false, ErrNotMergeCommit
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
