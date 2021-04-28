package github

import "github.com/shurcooL/githubv4"

// commitInfoQuery is the GraphQL query for retrieving GitHub CI/CD status info for a specific commit
type commitInfoQuery struct {
	Resource struct {
		Typename string `graphql:"__typename"`
		Commit   struct {
			URL string
			// External services can mark commits with a Status that is reflected in pull requests involving those commits
			Status struct {
				State statusState
			}
			AssociatedPullRequests struct {
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
				Nodes []pullRequest
			} `graphql:"associatedPullRequests(first: 100, after: $associatedPullRequestsCursor)"`
			// CheckSuites is a collection of the check runs created by a single GitHub App for a specific commit
			// More info on CheckSuites here: https://developer.github.com/v4/guides/intro-to-graphql/#discovering-the-graphql-api
			CheckSuites struct {
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
				Nodes []checkSuite
			} `graphql:"checkSuites(first: 100, after: $checkSuitesCursor)"`
			Signature struct {
				IsValid bool
			}
			Repository struct {
				URL string
			}
		} `graphql:"... on Commit"`
	} `graphql:"resource(url: $url)"`
}

// checkSuite is a collection of the check runs created by a CI/CD App
type checkSuite struct {
	Status     checkStatusState
	Conclusion checkConclusionState
}

// pullRequest contains the relevant information associated with a pull request
type pullRequest struct {
	Merged      bool
	MergeCommit commit
	BaseRefName string
	HeadRefName string
	URL         string
}
