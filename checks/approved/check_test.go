package approved

import (
	"context"
	"testing"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/repository"
	r "github.com/grafeas/voucher/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprovedCheck(t *testing.T) {
	ctx := context.Background()
	imageData, err := voucher.NewImageData("gcr.io/voucher-test-project/apps/staging/voucher-internal@sha256:73d506a23331fce5cb6f49bfb4c27450d2ef4878efce89f03a46b27372a88430")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)
	buildDetail := r.BuildDetail{RepositoryURL: "https://github.com/grafeas/voucher-internal", Commit: "efgh6543"}
	commitURL := "https://github.com/grafeas/voucher-internal/commit/efgh6543"

	cases := []struct {
		name                 string
		defaultBranchCommits []r.CommitRef
		isSigned             bool
		status               string
		pullRequest          r.PullRequest
		shouldPass           bool
		err                  error
	}{
		{
			name:                 "Should pass",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             true,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}, HasRequiredApprovals: true},
			shouldPass:           true,
			err:                  nil,
		},
		{
			name:                 "Not built off default branch",
			defaultBranchCommits: []r.CommitRef{{URL: "otherCommit"}},
			isSigned:             true,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}, HasRequiredApprovals: true},
			shouldPass:           false,
			err:                  ErrNotOnDefaultBranch,
		},
		{
			name:                 "Commit not signed",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             false,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}, HasRequiredApprovals: true},
			shouldPass:           false,
			err:                  ErrNotSigned,
		},
		{
			name:                 "Commit not a merge commit",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             true,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: "otherURL"}, HasRequiredApprovals: true},
			shouldPass:           false,
			err:                  ErrNotMergeCommit,
		},
		{
			name:                 "Commit PR does not have required approvals",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             true,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}, HasRequiredApprovals: false},
			shouldPass:           false,
			err:                  ErrMissingRequiredApprovals,
		},
		{
			name:                 "CI check not successful",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             true,
			status:               "FAILURE",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}, HasRequiredApprovals: true},
			shouldPass:           false,
			err:                  ErrNotPassedCI,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			commit := r.Commit{
				URL:                    commitURL,
				Status:                 testCase.status,
				IsSigned:               testCase.isSigned,
				AssociatedPullRequests: []r.PullRequest{testCase.pullRequest},
			}
			defaultBranch := r.Branch{Name: "production", CommitRefs: testCase.defaultBranchCommits}

			metadataClient := new(voucher.MockMetadataClient)
			metadataClient.On("GetBuildDetail", ctx, imageData).Return(buildDetail, nil)

			repositoryClient := new(repository.MockClient)
			repositoryClient.On("GetCommit", ctx, buildDetail).Return(commit, nil)
			repositoryClient.On("GetDefaultBranch", ctx, buildDetail).Return(defaultBranch, nil)

			orgCheck := new(check)
			orgCheck.SetMetadataClient(metadataClient)
			orgCheck.SetRepositoryClient(repositoryClient)

			status, err := orgCheck.Check(ctx, imageData)

			assert.Equal(t, testCase.shouldPass, status)
			if testCase.err != nil {
				assert.EqualError(t, testCase.err, err.Error())
			}
		})
	}
}
