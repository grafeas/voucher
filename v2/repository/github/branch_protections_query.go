package github

import "github.com/shurcooL/githubv4"

// branchProtectionsQuery is the GraphQL query for retrieving information pertaining to branch protections in a repository
type branchProtectionsQuery struct {
	Resource struct {
		Typename   string `graphql:"__typename"`
		Repository struct {
			BranchProtectionRules struct {
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
				Nodes []branchProtectionRule
			} `graphql:"branchProtectionRules(first:100, after: $branchProtectionRulesCursor)"`
		} `graphql:"... on Repository"`
	} `graphql:"resource(url: $url)"`
}

type branchProtectionRule struct {
	RequiresApprovingReviews     bool
	RequiredApprovingReviewCount int
	MatchingRefs                 struct {
		PageInfo struct {
			EndCursor   githubv4.String
			HasNextPage bool
		}
		Nodes []matchingRef
	} `graphql:"matchingRefs(first: 100, after: $matchingRefsCursor)"`
}

type matchingRef struct {
	Name string
}
