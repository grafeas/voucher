package github

import "github.com/shurcooL/githubv4"

// branchQuery is the GraphQL query for retrieving information pertaining to a branch in a repository
type branchQuery struct {
	Resource struct {
		Typename   string `graphql:"__typename"`
		Repository struct {
			Ref struct {
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
							Nodes    []commit // Nodes contains all of the commits in the branch
						} `graphql:"history(first: 100, after: $branchCommitCursor)"`
					} `graphql:"... on Commit"`
				}
			} `graphql:"ref(qualifiedName: $branch_name)"`
		} `graphql:"... on Repository"`
	} `graphql:"resource(url: $url)"`
}
