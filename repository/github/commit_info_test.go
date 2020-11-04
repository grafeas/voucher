package github

import (
	"context"
	"testing"

	"github.com/grafeas/voucher/repository"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllCheckSuites(t *testing.T) {
	getAllCheckSuitesTests := []struct {
		testName                 string
		commitURL                string
		input                    *commitInfoQuery
		mask                     []string
		queryPopulationVariables map[string]interface{}
		expected                 []checkSuite
	}{
		{
			testName:  "Testing zero check suites",
			commitURL: "https://github.com/grafeas/voucher/commit/881e9fd71e816415e1f199daeb6dc6d3c5fd4f2c",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/grafeas/voucher/commit/881e9fd71e816415e1f199daeb6dc6d3c5fd4f2c"
				res.Resource.Commit.CheckSuites.Nodes = []checkSuite{}
				res.Resource.Commit.CheckSuites.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Commit.CheckSuites.Nodes",
			},
			expected: []checkSuite{},
		},
		{
			testName:  "Testing happy check suites",
			commitURL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"
				res.Resource.Commit.CheckSuites.Nodes = []checkSuite{
					{
						Status:     "COMPLETED",
						Conclusion: "SUCCESS",
					},
				}
				res.Resource.Commit.CheckSuites.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Commit.CheckSuites.Nodes",
			},
			expected: []checkSuite{
				{
					Status:     "COMPLETED",
					Conclusion: "SUCCESS",
				},
			},
		},
	}
	for _, test := range getAllCheckSuitesTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.input, test.mask)
			formattedURI, err := createNewGitHubV4URI(test.commitURL)
			require.Equal(t, commitType, test.input.Resource.Typename)

			assert.NoError(t, err)

			res, err := getAllCheckSuites(context.Background(), c, test.input, githubv4.URI(*formattedURI))
			assert.NoError(t, err, "Getting all check suites failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}

func TestGetAllAssociatedPullRequests(t *testing.T) {
	getAllAssociatedPullRequestsTests := []struct {
		testName                 string
		commitURL                string
		input                    *commitInfoQuery
		mask                     []string
		queryPopulationVariables map[string]interface{}
		expected                 []pullRequest
	}{
		{
			testName:  "Testing zero associated pull requests",
			commitURL: "https://github.com/grafeas/voucher/commit/881e9fd71e816415e1f199daeb6dc6d3c5fd4f2c",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/grafeas/voucher/commit/881e9fd71e816415e1f199daeb6dc6d3c5fd4f2c"
				res.Resource.Commit.AssociatedPullRequests.Nodes = []pullRequest{}
				res.Resource.Commit.AssociatedPullRequests.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Commit.AssociatedPullRequests.Nodes",
			},
			expected: []pullRequest{},
		},
		{
			testName:  "Testing has associated pull requests",
			commitURL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"
				res.Resource.Commit.AssociatedPullRequests.Nodes = []pullRequest{
					{
						BaseRefName: "master",
						HeadRefName: "fix-broken-response",
						Merged:      true,
						MergeCommit: commit{
							URL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
						},
						URL: "https://github.com/grafeas/voucher/pull/23",
					},
				}
				res.Resource.Commit.AssociatedPullRequests.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Commit.AssociatedPullRequests.Nodes",
			},
			expected: []pullRequest{
				{
					BaseRefName: "master",
					HeadRefName: "fix-broken-response",
					Merged:      true,
					MergeCommit: commit{
						URL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
					},
					URL: "https://github.com/grafeas/voucher/pull/23",
				},
			},
		},
	}
	for _, test := range getAllAssociatedPullRequestsTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.input, test.mask)
			formattedURI, err := createNewGitHubV4URI(test.commitURL)
			require.Equal(t, commitType, test.input.Resource.Typename)

			assert.NoError(t, err)

			res, err := getAllAssociatedPullRequests(context.Background(), c, test.input, githubv4.URI(*formattedURI))
			assert.NoError(t, err, "Getting all associated pull requests failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}

func TestConvertCheckSuites(t *testing.T) {
	testCases := []struct {
		testName    string
		checkSuites []checkSuite
		expected    []repository.Check
	}{
		{
			testName:    "Testing zero check suites",
			checkSuites: []checkSuite{},
			expected:    []repository.Check{},
		},
		{
			testName: "Testing check suites",
			checkSuites: []checkSuite{
				{
					Status:     "COMPLETED",
					Conclusion: "SUCCESS",
				},
				{
					Status:     "COMPLETED",
					Conclusion: "FAILURE",
				},
			},
			expected: []repository.Check{
				{
					Status:     "COMPLETED",
					Conclusion: "SUCCESS",
				},
				{
					Status:     "COMPLETED",
					Conclusion: "FAILURE",
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			res, err := convertCheckSuites(test.checkSuites)
			assert.NoError(t, err, "Converting checkSuites to repository.Checks failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}

func TestConvertPullRequests(t *testing.T) {
	testCases := []struct {
		testName                 string
		pullRequestURL           string
		pullRequests             []pullRequest
		branchProtections        map[string]branchProtection
		input                    *pullRequestReviewsQuery
		mask                     []string
		queryPopulationVariables map[string]interface{}
		expected                 []repository.PullRequest
	}{
		{
			testName:     "Testing zero pull requests",
			pullRequests: []pullRequest{},
			input: func() *pullRequestReviewsQuery {
				res := new(pullRequestReviewsQuery)
				res.Resource.Typename = "PullRequest"
				return res
			}(),
			expected: []repository.PullRequest{},
		},
		{
			testName:       "Testing has required approvals",
			pullRequestURL: "https://github.com/grafeas/voucher/pull/23",
			pullRequests: []pullRequest{
				{
					BaseRefName: "master",
					HeadRefName: "fix-broken-response",
					Merged:      true,
					MergeCommit: commit{
						URL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
					},
					URL: "https://github.com/grafeas/voucher/pull/23",
				},
			},
			branchProtections: map[string]branchProtection{
				"master": {
					RequiresApprovingReviews:     true,
					RequiredApprovingReviewCount: 2,
				},
			},
			input: func() *pullRequestReviewsQuery {
				res := new(pullRequestReviewsQuery)
				res.Resource.Typename = "PullRequest"
				res.Resource.PullRequest.Reviews.Nodes = []review{
					{State: "APPROVED"},
					{State: "APPROVED"},
				}
				res.Resource.PullRequest.Reviews.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.PullRequest.Reviews.Nodes",
			},
			expected: []repository.PullRequest{
				{
					BaseBranchName: "master",
					HeadBranchName: "fix-broken-response",
					IsMerged:       true,
					MergeCommit: repository.CommitRef{
						URL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
					},
					HasRequiredApprovals: true,
				},
			},
		},
		{
			testName:       "Testing not enough approvals",
			pullRequestURL: "https://github.com/grafeas/voucher/pull/23",
			pullRequests: []pullRequest{
				{
					BaseRefName: "master",
					HeadRefName: "fix-broken-response",
					Merged:      true,
					MergeCommit: commit{
						URL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
					},
					URL: "https://github.com/grafeas/voucher/pull/23",
				},
			},
			branchProtections: map[string]branchProtection{
				"master": {
					RequiresApprovingReviews:     true,
					RequiredApprovingReviewCount: 2,
				},
			},
			input: func() *pullRequestReviewsQuery {
				res := new(pullRequestReviewsQuery)
				res.Resource.Typename = "PullRequest"
				res.Resource.PullRequest.Reviews.Nodes = []review{
					{State: "APPROVED"},
					{State: "PENDING"},
				}
				res.Resource.PullRequest.Reviews.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.PullRequest.Reviews.Nodes",
			},
			expected: []repository.PullRequest{
				{
					BaseBranchName: "master",
					HeadBranchName: "fix-broken-response",
					IsMerged:       true,
					MergeCommit: repository.CommitRef{
						URL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
					},
					HasRequiredApprovals: false,
				},
			},
		},
		{
			testName:       "Testing approvals not required",
			pullRequestURL: "https://github.com/grafeas/voucher/pull/23",
			pullRequests: []pullRequest{
				{
					BaseRefName: "master",
					HeadRefName: "fix-broken-response",
					Merged:      true,
					MergeCommit: commit{
						URL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
					},
					URL: "https://github.com/grafeas/voucher/pull/23",
				},
			},
			branchProtections: map[string]branchProtection{
				"master": {
					RequiresApprovingReviews:     false,
					RequiredApprovingReviewCount: 2,
				},
			},
			input: func() *pullRequestReviewsQuery {
				res := new(pullRequestReviewsQuery)
				res.Resource.Typename = "PullRequest"
				res.Resource.PullRequest.Reviews.Nodes = []review{}
				res.Resource.PullRequest.Reviews.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.PullRequest.Reviews.Nodes",
			},
			expected: []repository.PullRequest{
				{
					BaseBranchName: "master",
					HeadBranchName: "fix-broken-response",
					IsMerged:       true,
					MergeCommit: repository.CommitRef{
						URL: "https://github.com/grafeas/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
					},
					HasRequiredApprovals: true,
				},
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.input, test.mask)
			require.Equal(t, pullRequestType, test.input.Resource.Typename)

			res, err := convertPullRequests(context.Background(), c, test.pullRequests, test.branchProtections)
			assert.NoError(t, err, "Converting pullRequests to repository.PullRequests failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}
