package github

import (
	"context"

	"github.com/grafeas/voucher/repository"
	"github.com/shurcooL/githubv4"
)

// getAllReviews is the graphQL query for collecting all reviews associated with a pull request
func getAllReviews(ctx context.Context, ghc ghGraphQLClient, pullRequestURL string) ([]review, error) {
	allReviews := make([]review, 0)
	formattedURI, err := createNewGitHubV4URI(pullRequestURL)
	if err != nil {
		return nil, err
	}
	queryResult := new(pullRequestReviewsQuery)
	pullRequestVariables := map[string]interface{}{
		"url":           githubv4.URI(*formattedURI),
		"reviewsCursor": (*githubv4.String)(nil),
	}
	err = paginationQuery(ctx, ghc, queryResult, pullRequestVariables, queryPageLimit, func(v interface{}) (bool, error) {
		ciq, ok := v.(*pullRequestReviewsQuery)
		if !ok {
			return false, newTypeMismatchError("pullRequestReviewsQuery", ciq)
		}
		pullRequest := ciq.Resource.PullRequest
		resourceType := ciq.Resource.Typename
		if resourceType != pullRequestType {
			return false, repository.NewTypeMismatchError(commitType, resourceType)
		}

		allReviews = append(allReviews, pullRequest.Reviews.Nodes...)
		hasMoreResults := pullRequest.Reviews.PageInfo.HasNextPage
		pullRequestVariables["reviewsCursor"] = githubv4.NewString(pullRequest.Reviews.PageInfo.EndCursor)
		return hasMoreResults, nil
	})
	if err != nil {
		return nil, err
	}
	return allReviews, nil
}
