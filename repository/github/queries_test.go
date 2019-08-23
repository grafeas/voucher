package github

import (
	"context"
	"fmt"
	"testing"

	"github.com/Shopify/voucher/repository"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	utils "github.com/mennanov/fieldmask-utils"
	"github.com/shurcooL/githubv4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type queryHandlerFunc func(query interface{}, variables map[string]interface{}) error

func createHandler(commitURL string, input interface{}, mask []string, expected interface{}) queryHandlerFunc {
	return func(query interface{}, variables map[string]interface{}) error {
		mask, err := utils.MaskFromPaths(mask, generator.CamelCase)
		if err != nil {
			return err
		}

		err = utils.StructToStruct(mask, input, query)
		if err != nil {
			return err
		}

		return nil
	}
}

type mockGitHubGraphQLClient struct {
	HandlerFunc queryHandlerFunc
}

func (m *mockGitHubGraphQLClient) Query(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	return m.HandlerFunc(query, variables)
}

func TestNewRepositoryOrgInfoResult(t *testing.T) {
	newRepoOrgTests := []struct {
		testName  string
		commitURL string
		input     *repositoryOrgInfoQuery
		mask      []string
		expected  *repositoryOrgInfoQuery
	}{
		{
			testName:  "Testing happy path",
			commitURL: "www.shopify.com",
			input: func() *repositoryOrgInfoQuery {
				res := new(repositoryOrgInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.Repository.Owner.Typename = "Organization"
				org := res.Resource.Commit.Repository.Owner.Organization
				org.ID = "2342ffesfdf"
				org.Name = "Shopify"
				org.URL = "github.com/Shopify"
				return res
			}(),
			mask: []string{"Resource.Typename", "Resource.Commit.Repository.Owner.Typename"},
			expected: func() *repositoryOrgInfoQuery {
				res := new(repositoryOrgInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.Repository.Owner.Typename = "Organization"
				return res
			}(),
		},
		{
			testName:  "Testing with bad URL",
			commitURL: "hello@%a&%(.com",
			input:     new(repositoryOrgInfoQuery),
			mask:      []string{},
			expected:  nil,
		},
	}

	for _, test := range newRepoOrgTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.commitURL, test.input, test.mask, test.expected)
			res, err := newRepositoryOrgInfoResult(context.Background(), c, test.commitURL)

			if test.expected == nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Exactly(t, test.expected, res)

			require.Equal(t, commitType, res.Resource.Typename)
			assert.Equal(t, organizationType, res.Resource.Commit.Repository.Owner.Typename)
		})
	}
}

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
			commitURL: "https://github.com/Shopify/voucher/commit/881e9fd71e816415e1f199daeb6dc6d3c5fd4f2c",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/Shopify/voucher/commit/881e9fd71e816415e1f199daeb6dc6d3c5fd4f2c"
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
			commitURL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"
				res.Resource.Commit.CheckSuites.Nodes = []checkSuite{
					{
						App: app{
							Name: "Travis CI",
							URL:  "https://travis-ci.com",
						},
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
					App: app{
						Name: "Travis CI",
						URL:  "https://travis-ci.com",
					},
					Status:     "COMPLETED",
					Conclusion: "SUCCESS",
				},
			},
		},
	}
	for _, test := range getAllCheckSuitesTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.commitURL, test.input, test.mask, test.expected)
			formattedURI, err := createNewGitHubV4URI(test.commitURL)
			queryPopulationVariables := map[string]interface{}{
				"url":               githubv4.URI(*formattedURI),
				"checkSuitesCursor": (*githubv4.String)(nil),
			}
			require.Equal(t, commitType, test.input.Resource.Typename)

			assert.NoError(t, err)

			res, err := getAllCheckSuites(context.Background(), c, test.input, queryPopulationVariables)
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
			commitURL: "https://github.com/Shopify/voucher/commit/881e9fd71e816415e1f199daeb6dc6d3c5fd4f2c",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/Shopify/voucher/commit/881e9fd71e816415e1f199daeb6dc6d3c5fd4f2c"
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
			commitURL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"
				res.Resource.Commit.AssociatedPullRequests.Nodes = []pullRequest{
					{
						BaseRefName: "master",
						HeadRefName: "fix-broken-response",
						Merged:      true,
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
				},
			},
		},
	}
	for _, test := range getAllAssociatedPullRequestsTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.commitURL, test.input, test.mask, test.expected)
			formattedURI, err := createNewGitHubV4URI(test.commitURL)
			queryPopulationVariables := map[string]interface{}{
				"url":                          githubv4.URI(*formattedURI),
				"associatedPullRequestsCursor": (*githubv4.String)(nil),
			}
			require.Equal(t, commitType, test.input.Resource.Typename)

			assert.NoError(t, err)

			res, err := getAllAssociatedPullRequests(context.Background(), c, test.input, queryPopulationVariables)
			assert.NoError(t, err, "Getting all associated pull requests failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}

func TestPaginationQuery(t *testing.T) {
	paginationTests := []struct {
		testName                 string
		queryResult              interface{}
		queryPopulationVariables map[string]interface{}
		pageLimit                int
		resultHandler            queryHandler
		errorExpected            error
	}{
		{
			testName:                 "Test infinite result with finite page limit",
			queryResult:              struct{}{},
			queryPopulationVariables: map[string]interface{}{},
			pageLimit:                2,
			resultHandler: func(queryResult interface{}) (bool, error) {
				return true, nil
			},
			errorExpected: nil,
		},
		{
			testName:                 "Test error result",
			queryResult:              struct{}{},
			queryPopulationVariables: map[string]interface{}{},
			pageLimit:                2,
			resultHandler: func(queryResult interface{}) (bool, error) {
				return true, fmt.Errorf("cannot paginate anymore")
			},
			errorExpected: fmt.Errorf("cannot paginate anymore"),
		},
		{
			testName:                 "Test false result",
			queryResult:              struct{}{},
			queryPopulationVariables: map[string]interface{}{},
			pageLimit:                2,
			resultHandler: func(queryResult interface{}) (bool, error) {
				return false, nil
			},
			errorExpected: nil,
		},
	}

	for _, test := range paginationTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = func(query interface{}, variables map[string]interface{}) error {
				return nil
			}

			err := paginationQuery(
				context.Background(),
				c,
				test.queryResult,
				test.queryPopulationVariables,
				test.pageLimit,
				test.resultHandler,
			)

			if test.errorExpected != nil {
				assert.Equal(t, test.errorExpected, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestNewCommitInfo(t *testing.T) {
	commitInfoTests := []struct {
		testName               string
		queryResult            *commitInfoQuery
		checkSuites            []checkSuite
		associatedPullRequests []pullRequest
		expected               repository.CommitInfo
		shouldError            bool
	}{
		{
			testName: "Testing a healthy commit",
			queryResult: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"
				res.Resource.Commit.CheckSuites.Nodes = []checkSuite{
					{
						App: app{
							Name: "Travis CI",
							URL:  "https://travis-ci.com",
						},
						Status:     "COMPLETED",
						Conclusion: "SUCCESS",
					},
				}
				res.Resource.Commit.CheckSuites.PageInfo.HasNextPage = false
				res.Resource.Commit.AssociatedPullRequests.Nodes = []pullRequest{
					{
						BaseRefName: "master",
						HeadRefName: "fix-broken-response",
						Merged:      true,
					},
				}
				res.Resource.Commit.Signature.IsValid = true
				res.Resource.Commit.Status.State = "SUCCESS"
				return res
			}(),
			checkSuites: []checkSuite{
				{
					App: app{
						Name: "Travis CI",
						URL:  "https://travis-ci.com",
					},
					Status:     "COMPLETED",
					Conclusion: "SUCCESS",
				},
			},
			associatedPullRequests: []pullRequest{
				{
					BaseRefName: "master",
					HeadRefName: "fix-broken-response",
					Merged:      true,
				},
			},
			expected: repository.CommitInfo{
				URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
				Checks: []repository.Check{
					{
						App: repository.App{
							Name: "Travis CI",
							URL:  "https://travis-ci.com",
						},
						Status:     "COMPLETED",
						Conclusion: "SUCCESS",
					},
				},
				Status:   "SUCCESS",
				IsSigned: true,
				AssociatedPullRequests: []repository.PullRequest{
					{
						BaseBranchName: "master",
						HeadBranchName: "fix-broken-response",
						IsMerged:       true,
					},
				},
			},
		},
		{
			testName: "Testing invalid state",
			queryResult: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Commit.Status.State = "invalid"
				return res
			}(),
			checkSuites:            []checkSuite{},
			associatedPullRequests: []pullRequest{},
			expected:               repository.CommitInfo{},
			shouldError:            true,
		},
	}
	for _, test := range commitInfoTests {
		t.Run(test.testName, func(t *testing.T) {
			res, err := createNewCommitInfo(test.queryResult, test.checkSuites, test.associatedPullRequests)
			if test.shouldError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err, "Creating new CommitInfo object failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}

func TestNewDefaultBranchResult(t *testing.T) {
	defaultBranchResultTests := []struct {
		testName    string
		commitURL   string
		input       *defaultBranchQuery
		mask        []string
		expected    repository.DefaultBranch
		shouldError bool
	}{
		{
			testName:  "Testing no commits in default branch",
			commitURL: "github.com/Shopify/voucher/commit/FakeCommit",
			input: func() *defaultBranchQuery {
				res := new(defaultBranchQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.Repository.DefaultBranchRef.Name = "master"
				res.Resource.Commit.Repository.DefaultBranchRef.Target.Commit.Typename = "Commit"
				res.Resource.Commit.Repository.DefaultBranchRef.Target.Commit.History.Nodes = []commit{}
				return res
			}(),
			mask: []string{"Resource.Typename", "Resource.Commit.Repository.DefaultBranchRef.Target.Commit.History.Nodes"},
			expected: repository.DefaultBranch{
				Name:    "",
				Commits: []repository.Commit{},
			},
			shouldError: false,
		},
		{
			testName:  "Testing has some commits in default branch",
			commitURL: "github.com/Shopify/voucher/commit/FakeCommit1",
			input: func() *defaultBranchQuery {
				res := new(defaultBranchQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.Repository.DefaultBranchRef.Name = "master"
				res.Resource.Commit.Repository.DefaultBranchRef.Target.Commit.Typename = "Commit"
				res.Resource.Commit.Repository.DefaultBranchRef.Target.Commit.History.Nodes = []commit{
					{
						URL: "github.com/Shopify/voucher/commit/FakeCommit1",
					},
					{
						URL: "github.com/Shopify/voucher/commit/FakeCommit2",
					},
				}
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Commit.Repository.DefaultBranchRef",
			},
			expected: repository.DefaultBranch{
				Name: "master",
				Commits: []repository.Commit{
					{
						URL: "github.com/Shopify/voucher/commit/FakeCommit1",
					},
					{
						URL: "github.com/Shopify/voucher/commit/FakeCommit2",
					},
				},
			},
			shouldError: false,
		},
		{
			testName:    "Testing error propagation",
			commitURL:   "github.com/Shopify/voucher/commit/%$#^&@",
			input:       new(defaultBranchQuery),
			mask:        []string{},
			expected:    repository.DefaultBranch{},
			shouldError: true,
		},
	}
	for _, test := range defaultBranchResultTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.commitURL, test.input, test.mask, test.expected)
			res, err := newDefaultBranchResult(context.Background(), c, test.commitURL)
			if test.shouldError {
				assert.Error(t, err)
				return
			}
			require.Equal(t, commitType, test.input.Resource.Typename)

			assert.NoError(t, err, "Getting default branch result failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}

func TestCreateNewDefaultBranch(t *testing.T) {
	createNewDefaultBranchTests := []struct {
		testName               string
		queryResult            *commitInfoQuery
		checkSuites            []checkSuite
		associatedPullRequests []pullRequest
		expected               repository.CommitInfo
	}{
		{
			testName: "Testing an unsigned commit",
			queryResult: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"
				res.Resource.Commit.CheckSuites.Nodes = []checkSuite{
					{
						App: app{
							Name: "Travis CI",
							URL:  "https://travis-ci.com",
						},
						Status:     "COMPLETED",
						Conclusion: "SUCCESS",
					},
				}
				res.Resource.Commit.CheckSuites.PageInfo.HasNextPage = false
				res.Resource.Commit.AssociatedPullRequests.Nodes = []pullRequest{
					{
						BaseRefName: "master",
						HeadRefName: "fix-broken-response",
						Merged:      true,
					},
				}
				res.Resource.Commit.Signature.IsValid = true
				res.Resource.Commit.Status.State = "SUCCESS"
				return res
			}(),
			checkSuites: []checkSuite{
				{
					App: app{
						Name: "Travis CI",
						URL:  "https://travis-ci.com",
					},
					Status:     "COMPLETED",
					Conclusion: "SUCCESS",
				},
			},
			associatedPullRequests: []pullRequest{
				{
					BaseRefName: "master",
					HeadRefName: "fix-broken-response",
					Merged:      true,
				},
			},
			expected: repository.CommitInfo{
				URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
				Checks: []repository.Check{
					{
						App: repository.App{
							Name: "Travis CI",
							URL:  "https://travis-ci.com",
						},
						Status:     "COMPLETED",
						Conclusion: "SUCCESS",
					},
				},
				Status:   "SUCCESS",
				IsSigned: true,
				AssociatedPullRequests: []repository.PullRequest{
					{
						BaseBranchName: "master",
						HeadBranchName: "fix-broken-response",
						IsMerged:       true,
					},
				},
			},
		},
	}
	for _, test := range createNewDefaultBranchTests {
		t.Run(test.testName, func(t *testing.T) {
			res, err := createNewCommitInfo(test.queryResult, test.checkSuites, test.associatedPullRequests)
			require.NoError(t, err, "Creating new CommitInfo object failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}

func TestNewCommitInfoResult(t *testing.T) {
	newCommitInfoResultTests := []struct {
		testName    string
		commitURL   string
		input       *commitInfoQuery
		mask        []string
		expected    repository.CommitInfo
		shouldError bool
	}{
		{
			testName:  "Testing bad commit",
			commitURL: "https://github.com/Shopify/voucher/&%$%^&)",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/Shopify/voucher/&%$%^&)"
				res.Resource.Commit.AssociatedPullRequests.Nodes = []pullRequest{}
				res.Resource.Commit.AssociatedPullRequests.PageInfo.HasNextPage = false
				res.Resource.Commit.CheckSuites.Nodes = []checkSuite{}
				res.Resource.Commit.CheckSuites.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Commit.AssociatedPullRequests.Nodes",
				"Resource.Commit.CheckSuites.Nodes",
			},
			expected:    repository.CommitInfo{},
			shouldError: true,
		},
		{
			testName:  "Testing has check suites and associated pull requests",
			commitURL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
			input: func() *commitInfoQuery {
				res := new(commitInfoQuery)
				res.Resource.Typename = "Commit"
				res.Resource.Commit.URL = "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"
				res.Resource.Commit.AssociatedPullRequests.Nodes = []pullRequest{
					{
						BaseRefName: "master",
						HeadRefName: "fix-broken-response",
						Merged:      true,
					},
				}
				res.Resource.Commit.Signature.IsValid = true
				res.Resource.Commit.Status.State = "SUCCESS"
				res.Resource.Commit.AssociatedPullRequests.PageInfo.HasNextPage = false
				res.Resource.Commit.CheckSuites.Nodes = []checkSuite{
					{
						App: app{
							Name: "Travis CI",
							URL:  "https://travis-ci.com",
						},
						Status:     "COMPLETED",
						Conclusion: "SUCCESS",
					},
				}
				res.Resource.Commit.CheckSuites.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Commit",
			},
			expected: repository.CommitInfo{
				URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
				Checks: []repository.Check{
					{
						App: repository.App{
							Name: "Travis CI",
							URL:  "https://travis-ci.com",
						},
						Status:     "COMPLETED",
						Conclusion: "SUCCESS",
					},
				},
				Status:   "SUCCESS",
				IsSigned: true,
				AssociatedPullRequests: []repository.PullRequest{
					{
						BaseBranchName: "master",
						HeadBranchName: "fix-broken-response",
						IsMerged:       true,
					},
				},
			},
			shouldError: false,
		},
	}
	for _, test := range newCommitInfoResultTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.commitURL, test.input, test.mask, test.expected)

			res, err := newCommitInfoResult(context.Background(), c, test.commitURL)

			if test.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err, "Getting all commit information failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}
