package github

import (
	"context"
	"testing"

	"github.com/grafeas/voucher/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBranchResult(t *testing.T) {
	branchResultTests := []struct {
		testName    string
		repoURL     string
		input       *branchQuery
		mask        []string
		expected    repository.Branch
		shouldError bool
	}{
		{
			testName: "Testing no commits in master branch",
			repoURL:  "github.com/grafeas/voucher",
			input: func() *branchQuery {
				res := new(branchQuery)
				res.Resource.Typename = "Repository"
				res.Resource.Repository.Ref.Name = "master"
				res.Resource.Repository.Ref.Target.Commit.Typename = "Commit"
				res.Resource.Repository.Ref.Target.Commit.History.Nodes = []commit{}
				return res
			}(),
			mask: []string{"Resource.Typename", "Resource.Repository.Ref"},
			expected: repository.Branch{
				Name:       "master",
				CommitRefs: []repository.CommitRef{},
			},
			shouldError: false,
		},
		{
			testName: "Testing has some commits in master branch",
			repoURL:  "github.com/grafeas/voucher",
			input: func() *branchQuery {
				res := new(branchQuery)
				res.Resource.Typename = "Repository"
				res.Resource.Repository.Ref.Name = "master"
				res.Resource.Repository.Ref.Target.Commit.Typename = "Commit"
				res.Resource.Repository.Ref.Target.Commit.History.Nodes = []commit{
					{
						URL: "github.com/grafeas/voucher/commit/FakeCommit1",
					},
					{
						URL: "github.com/grafeas/voucher/commit/FakeCommit2",
					},
				}
				return res
			}(),
			mask: []string{
				"Resource.Typename", "Resource.Repository.Ref",
			},
			expected: repository.Branch{
				Name: "master",
				CommitRefs: []repository.CommitRef{
					{
						URL: "github.com/grafeas/voucher/commit/FakeCommit1",
					},
					{
						URL: "github.com/grafeas/voucher/commit/FakeCommit2",
					},
				},
			},
			shouldError: false,
		},
		{
			testName:    "Testing error propagation",
			repoURL:     "error.com/grafeas/voucher",
			input:       new(branchQuery),
			mask:        []string{},
			expected:    repository.Branch{},
			shouldError: true,
		},
	}
	for _, test := range branchResultTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.input, test.mask)
			res, err := newBranchResult(context.Background(), c, test.repoURL, "master")
			if test.shouldError {
				assert.Error(t, err)
				return
			}
			require.Equal(t, repositoryType, test.input.Resource.Typename)

			assert.NoError(t, err, "Getting master branch result failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}
