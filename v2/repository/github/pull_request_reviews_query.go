package github

import "github.com/shurcooL/githubv4"

// pullRequestReviewsQuery is the GraphQL query for retrieving information pertaining to pull request reviews
type pullRequestReviewsQuery struct {
	Resource struct {
		Typename    string `graphql:"__typename"`
		PullRequest struct {
			Reviews struct {
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
				Nodes []review
			} `graphql:"reviews(first: 100, after: $reviewsCursor)"`
		} `graphql:"... on PullRequest"`
	} `graphql:"resource(url: $url)"`
}

type review struct {
	State pullRequestReviewState
}
