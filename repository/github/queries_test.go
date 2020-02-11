package github

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/generator"
	utils "github.com/mennanov/fieldmask-utils"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Shopify/voucher/repository"
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
		testName    string
		uri         string
		input       *repositoryOrgInfoQuery
		mask        []string
		expected    repository.Organization
		shouldError bool
	}{
		{
			testName: "Testing happy path",
			uri:      "https://github.com/Shopify/voucher",
			input: func() *repositoryOrgInfoQuery {
				res := new(repositoryOrgInfoQuery)
				res.Resource.Repository.Owner.Typename = "Organization"
				res.Resource.Repository.Owner.Organization.Name = "Shopify"
				res.Resource.Repository.Owner.Organization.URL = "https://github.com/Shopify"
				return res
			}(),
			mask: []string{"Resource.Repository.Owner.Typename", "Resource.Repository.Owner.Organization"},
			expected: repository.Organization{
				Alias: "Shopify",
				VCS:   "github.com",
				Name:  "Shopify",
			},
			shouldError: false,
		},
		{
			testName:    "Testing with bad URL",
			uri:         "hello@%a&%(.com",
			input:       new(repositoryOrgInfoQuery),
			mask:        []string{},
			expected:    repository.Organization{},
			shouldError: true,
		},
	}

	for _, test := range newRepoOrgTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.uri, test.input, test.mask, test.expected)
			res, err := newRepositoryOrgInfoResult(context.Background(), c, test.uri)
			if test.shouldError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, test.expected, res)
			assert.Equal(t, organizationType, test.input.Resource.Repository.Owner.Typename)
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
			require.Equal(t, commitType, test.input.Resource.Typename)

			assert.NoError(t, err)

			res, err := getAllAssociatedPullRequests(context.Background(), c, test.input, githubv4.URI(*formattedURI))
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
		expected               repository.Commit
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
						MergeCommit: commit{URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"},
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
					MergeCommit: commit{URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"},
				},
			},
			expected: repository.Commit{
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
						MergeCommit: repository.CommitRef{
							URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
						},
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
			expected:               repository.Commit{},
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
		repoURL     string
		input       *defaultBranchQuery
		mask        []string
		expected    repository.Branch
		shouldError bool
	}{
		{
			testName: "Testing no commits in default branch",
			repoURL:  "github.com/Shopify/voucher",
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
			repoURL:  "github.com/Shopify/voucher",
			input: func() *defaultBranchQuery {
				res := new(defaultBranchQuery)
				res.Resource.Typename = "Repository"
				res.Resource.Repository.DefaultBranchRef.Name = "master"
				res.Resource.Repository.DefaultBranchRef.Target.Commit.Typename = "Commit"
				res.Resource.Repository.DefaultBranchRef.Target.Commit.History.Nodes = []commit{
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
				"Resource.Typename", "Resource.Repository.DefaultBranchRef",
			},
			expected: repository.Branch{
				Name: "master",
				CommitRefs: []repository.CommitRef{
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
			repoURL:     "random.com/Shopify/voucher",
			input:       new(defaultBranchQuery),
			mask:        []string{},
			expected:    repository.Branch{},
			shouldError: true,
		},
	}
	for _, test := range defaultBranchResultTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.repoURL, test.input, test.mask, test.expected)
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

func TestCreateNewDefaultBranch(t *testing.T) {
	createNewDefaultBranchTests := []struct {
		testName               string
		queryResult            *commitInfoQuery
		checkSuites            []checkSuite
		associatedPullRequests []pullRequest
		expected               repository.Commit
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
						MergeCommit: commit{URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"},
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
					MergeCommit: commit{
						URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
					},
				},
			},
			expected: repository.Commit{
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
						MergeCommit: repository.CommitRef{
							URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
						},
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
			repoURL:  "github.com/Shopify/voucher",
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
			repoURL:  "github.com/Shopify/voucher",
			input: func() *branchQuery {
				res := new(branchQuery)
				res.Resource.Typename = "Repository"
				res.Resource.Repository.Ref.Name = "master"
				res.Resource.Repository.Ref.Target.Commit.Typename = "Commit"
				res.Resource.Repository.Ref.Target.Commit.History.Nodes = []commit{
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
				"Resource.Typename", "Resource.Repository.Ref",
			},
			expected: repository.Branch{
				Name: "master",
				CommitRefs: []repository.CommitRef{
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
			repoURL:     "error.com/Shopify/voucher",
			input:       new(branchQuery),
			mask:        []string{},
			expected:    repository.Branch{},
			shouldError: true,
		},
	}
	for _, test := range branchResultTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.repoURL, test.input, test.mask, test.expected)
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

func TestNewCommitInfoResult(t *testing.T) {
	newCommitInfoResultTests := []struct {
		testName    string
		commitURL   string
		input       *commitInfoQuery
		mask        []string
		expected    repository.Commit
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
			expected:    repository.Commit{},
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
						MergeCommit: commit{URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428"},
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
			expected: repository.Commit{
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
						MergeCommit: repository.CommitRef{
							URL: "https://github.com/Shopify/voucher/commit/8c235f3bd57393c53037b032e6da3e2b48aa0428",
						},
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
