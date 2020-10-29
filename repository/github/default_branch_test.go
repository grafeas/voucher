package github

import (
	"context"
	"testing"

	"github.com/grafeas/voucher/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultBranchResult(t *testing.T) {
	defaultBranchResultTests := []struct {
		testName    string
		repoURL     string
		input       *defaultBranchQuery
		mask        []string
		expected    repository.Branch
		shouldError bool
	}{
		{
			testName: "Testing no commits in default branch",
			repoURL:  "github.com/grafeas/voucher",
			input: func() *defaultBranchQuery {
				res := new(defaultBranchQuery)
				res.Resource.Typename = "Repository"
				res.Resource.Repository.DefaultBranchRef.Name = "master"
				res.Resource.Repository.DefaultBranchRef.Target.Commit.Typename = "Commit"
				res.Resource.Repository.DefaultBranchRef.Target.Commit.History.Nodes = []commit{}
				return res
			}(),
			mask: []string{"Resource.Typename", "Resource.Repository.DefaultBranchRef"},
			expected: repository.Branch{
				Name:       "master",
				CommitRefs: []repository.CommitRef{},
			},
			shouldError: false,
		},
		{
			testName: "Testing has some commits in default branch",
			repoURL:  "github.com/grafeas/voucher",
			input: func() *defaultBranchQuery {
				res := new(defaultBranchQuery)
				res.Resource.Typename = "Repository"
				res.Resource.Repository.DefaultBranchRef.Name = "master"
				res.Resource.Repository.DefaultBranchRef.Target.Commit.Typename = "Commit"
				res.Resource.Repository.DefaultBranchRef.Target.Commit.History.Nodes = []commit{
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
				"Resource.Typename", "Resource.Repository.DefaultBranchRef",
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
			repoURL:     "random.com/grafeas/voucher",
			input:       new(defaultBranchQuery),
			mask:        []string{},
			expected:    repository.Branch{},
			shouldError: true,
		},
	}
	for _, test := range defaultBranchResultTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.input, test.mask)
			res, err := newDefaultBranchResult(context.Background(), c, test.repoURL)
			if test.shouldError {
				assert.Error(t, err)
				return
			}
			require.Equal(t, repositoryType, test.input.Resource.Typename)

			assert.NoError(t, err, "Getting default branch result failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}
