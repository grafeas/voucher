package github

import "github.com/shurcooL/githubv4"

// defaultBranchQuery is the GraphQL query for retrieving information pertaining to the repository's default branch
type defaultBranchQuery struct {
	Resource struct {
		Typename   string `graphql:"__typename"`
		Repository struct {
			DefaultBranchRef struct {
				Name   string
				Target struct {
					Commit struct {
						Typename string `graphql:"__typename"`
						History  struct {
							PageInfo struct {
								EndCursor   githubv4.String
								HasNextPage bool
							}
							Typename string   `graphql:"__typename"`
							Nodes    []commit // Nodes contains all of the commits in the default branch
						} `graphql:"history(first: 100, after: $defaultBranchCommitCursor)"`
					} `graphql:"... on Commit"`
				}
			}
		} `graphql:"... on Repository"`
	} `graphql:"resource(url: $url)"`
}
