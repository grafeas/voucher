package github

import (
	"context"

	"github.com/shurcooL/githubv4"
)

// getBranchProtections is the graphQL query for collecting all branch protections for a repository
func getBranchProtections(ctx context.Context, ghc ghGraphQLClient, repoURL string) (map[string]branchProtection, error) {
	formattedURI, err := createNewGitHubV4URI(repoURL)
	if err != nil {
		return nil, err
	}
	queryResult := new(branchProtectionsQuery)
	branchProtectionVariables := map[string]interface{}{
		"url":                         githubv4.URI(*formattedURI),
		"branchProtectionRulesCursor": (*githubv4.String)(nil),
		"matchingRefsCursor":          (*githubv4.String)(nil),
	}

	branchProtections := make(map[string]branchProtection)

	err = paginationQuery(ctx, ghc, queryResult, branchProtectionVariables, queryPageLimit, func(v interface{}) (bool, error) {
		ciq, ok := v.(*branchProtectionsQuery)
		if !ok {
			return false, newTypeMismatchError("branchProtectionsQuery", ciq)
		}

		repository := ciq.Resource.Repository
		resourceType := ciq.Resource.Typename
		if resourceType != repositoryType {
			return false, newTypeMismatchError(repositoryType, resourceType)
		}

		for i := range repository.BranchProtectionRules.Nodes {
			err = paginationQuery(ctx, ghc, queryResult, branchProtectionVariables, queryPageLimit, func(v interface{}) (bool, error) {
				ciq, ok := v.(*branchProtectionsQuery)
				if !ok {
					return false, newTypeMismatchError("branchProtectionsQuery", ciq)
				}

				branchProtectionRulesNode := ciq.Resource.Repository.BranchProtectionRules.Nodes[i]

				for _, ref := range branchProtectionRulesNode.MatchingRefs.Nodes {
					branchProtections[ref.Name] = branchProtection{
						RequiresApprovingReviews:     branchProtectionRulesNode.RequiresApprovingReviews,
						RequiredApprovingReviewCount: branchProtectionRulesNode.RequiredApprovingReviewCount,
					}
				}

				hasMoreResults := branchProtectionRulesNode.MatchingRefs.PageInfo.HasNextPage
				branchProtectionVariables["matchingRefsCursor"] = githubv4.NewString(branchProtectionRulesNode.MatchingRefs.PageInfo.EndCursor)
				return hasMoreResults, nil
			})
			if err != nil {
				return false, err
			}
		}

		hasMoreResults := repository.BranchProtectionRules.PageInfo.HasNextPage
		branchProtectionVariables["branchProtectionRulesCursor"] = githubv4.NewString(repository.BranchProtectionRules.PageInfo.EndCursor)
		return hasMoreResults, nil
	})
	if err != nil {
		return nil, err
	}

	return branchProtections, nil
}

type branchProtection struct {
	RequiresApprovingReviews     bool
	RequiredApprovingReviewCount int
}
