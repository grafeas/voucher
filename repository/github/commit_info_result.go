package github

import (
	"context"

	"github.com/grafeas/voucher/repository"
	"github.com/shurcooL/githubv4"
)

// newCommitInfoResult calls the commitInfoQuery and populates the respective variables
func newCommitInfoResult(ctx context.Context, ghc ghGraphQLClient, commitURL string) (repository.Commit, error) {
	formattedURI, err := createNewGitHubV4URI(commitURL)
	if err != nil {
		return repository.Commit{}, err
	}
	queryResult := new(commitInfoQuery)

	associatedPullRequests, err := getAllAssociatedPullRequests(ctx, ghc, queryResult, githubv4.URI(*formattedURI))
	if err != nil {
		return repository.Commit{}, err
	}

	checkSuites, err := getAllCheckSuites(ctx, ghc, queryResult, githubv4.URI(*formattedURI))
	if err != nil {
		return repository.Commit{}, err
	}

	commit := queryResult.Resource.Commit
	status := commit.Status.State
	if !status.isValidStatusState() {
		return repository.Commit{}, newTypeMismatchError("statusState", status)
	}

	branchProtections, err := getBranchProtections(ctx, ghc, commit.Repository.URL)
	if err != nil {
		return repository.Commit{}, err
	}

	pullRequests, err := convertPullRequests(ctx, ghc, associatedPullRequests, branchProtections)
	if err != nil {
		return repository.Commit{}, err
	}

	checks, err := convertCheckSuites(checkSuites)
	if err != nil {
		return repository.Commit{}, err
	}

	isSigned := commit.Signature.IsValid

	return repository.NewCommit(commit.URL, checks, string(status), isSigned, pullRequests), nil
}

// convertPullRequests creates new PullRequest object for each pullRequest
func convertPullRequests(ctx context.Context, ghc ghGraphQLClient, associatedPullRequests []pullRequest, branchProtections map[string]branchProtection) ([]repository.PullRequest, error) {
	pullRequests := make([]repository.PullRequest, 0)
	for _, pr := range associatedPullRequests {
		commit := repository.NewCommitRef(pr.MergeCommit.URL)

		protections, ok := branchProtections[pr.BaseRefName]
		if !ok || !protections.RequiresApprovingReviews {
			pullRequests = append(pullRequests, repository.NewPullRequest(pr.BaseRefName, pr.HeadRefName, pr.Merged, commit, true))
			continue
		}

		reviews, err := getAllReviews(ctx, ghc, pr.URL)
		if err != nil {
			return nil, err
		}

		approvedReviews := 0
		for _, review := range reviews {
			if !review.State.isValidPullRequestReviewState() {
				return nil, newTypeMismatchError("pullRequestReviewState", review.State)
			}
			if review.State == pullRequestReviewApproved {
				approvedReviews++
			}
		}
		hasRequiredApprovals := approvedReviews >= protections.RequiredApprovingReviewCount
		pullRequests = append(pullRequests, repository.NewPullRequest(pr.BaseRefName, pr.HeadRefName, pr.Merged, commit, hasRequiredApprovals))
	}
	return pullRequests, nil
}

// convertCheckSuites creates a new repository.Check object for each checkSuite
func convertCheckSuites(checkSuites []checkSuite) ([]repository.Check, error) {
	checks := make([]repository.Check, 0)
	for _, checkSuite := range checkSuites {
		if !checkSuite.Status.isValidCheckStatusState() {
			return nil, newTypeMismatchError("checkStatusState", checkSuite.Status)
		}
		if !checkSuite.Conclusion.isValidCheckConclusionState() {
			return nil, newTypeMismatchError("checkConclusionState", checkSuite.Conclusion)
		}
		check := repository.NewCheck(string(checkSuite.Status), string(checkSuite.Conclusion))
		checks = append(checks, check)
	}
	return checks, nil
}

// getAllCheckSuites is the GraphQL query for collecting all the Check Suites pertaining to a commit
func getAllCheckSuites(ctx context.Context, ghc ghGraphQLClient, queryResult *commitInfoQuery, uri githubv4.URI) ([]checkSuite, error) {
	allCheckSuites := make([]checkSuite, 0)
	commitInfoVariables := map[string]interface{}{
		"url":                          uri,
		"checkSuitesCursor":            (*githubv4.String)(nil),
		"associatedPullRequestsCursor": (*githubv4.String)(nil),
	}
	err := paginationQuery(ctx, ghc, queryResult, commitInfoVariables, queryPageLimit, func(v interface{}) (bool, error) {
		ciq, ok := v.(*commitInfoQuery)
		if !ok {
			return false, newTypeMismatchError("commitInfoQuery", ciq)
		}
		commit := ciq.Resource.Commit
		resourceType := ciq.Resource.Typename
		if resourceType != commitType {
			return false, repository.NewTypeMismatchError(commitType, resourceType)
		}

		allCheckSuites = append(allCheckSuites, commit.CheckSuites.Nodes...)
		hasMoreResults := commit.CheckSuites.PageInfo.HasNextPage
		commitInfoVariables["checkSuitesCursor"] = githubv4.NewString(commit.CheckSuites.PageInfo.EndCursor)
		return hasMoreResults, nil
	})
	if err != nil {
		return nil, err
	}
	return allCheckSuites, nil
}

// getAllAssociatedPullRequests is the GraphQL query for collecting all the pull requests associated with a commit
func getAllAssociatedPullRequests(ctx context.Context, ghc ghGraphQLClient, queryResult *commitInfoQuery, uri githubv4.URI) ([]pullRequest, error) {
	allAssociatedPullRequests := make([]pullRequest, 0)
	commitInfoVariables := map[string]interface{}{
		"url":                          uri,
		"checkSuitesCursor":            (*githubv4.String)(nil),
		"associatedPullRequestsCursor": (*githubv4.String)(nil),
	}
	err := paginationQuery(ctx, ghc, queryResult, commitInfoVariables, queryPageLimit, func(v interface{}) (bool, error) {
		ciq, ok := v.(*commitInfoQuery)
		if !ok {
			return false, newTypeMismatchError("commitInfoQuery", ciq)
		}
		commit := ciq.Resource.Commit
		resourceType := ciq.Resource.Typename
		if resourceType != commitType {
			return false, repository.NewTypeMismatchError(commitType, resourceType)
		}

		allAssociatedPullRequests = append(allAssociatedPullRequests, commit.AssociatedPullRequests.Nodes...)
		hasMoreResults := commit.AssociatedPullRequests.PageInfo.HasNextPage
		commitInfoVariables["associatedPullRequestsCursor"] = githubv4.NewString(commit.AssociatedPullRequests.PageInfo.EndCursor)
		return hasMoreResults, nil
	})
	if err != nil {
		return nil, err
	}
	return allAssociatedPullRequests, nil
}
