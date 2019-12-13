package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"

	"github.com/Shopify/voucher/repository"
)

// queryHandler is called on every iteration of paginationQuery to populate a slice of query results
// queryHandler checks to see whether there are more records given that GitHub has a limit of 100 records per query
type queryHandler func(queryResult interface{}) (bool, error)

// paginationQuery populates a destination slice with the appropriately typed query results
// GitHub has a limit of 100 records so we must perform pagination
func paginationQuery(
	ctx context.Context,
	ghc ghGraphQLClient,
	queryResult interface{},
	queryPopulationVariables map[string]interface{},
	pageLimit int,
	qh queryHandler,
) error {
	for i := 0; i < pageLimit; i++ {
		err := ghc.Query(ctx, queryResult, queryPopulationVariables)
		if err != nil {
			return err
		}

		hasMoreResults, err := qh(queryResult)
		if nil != err {
			return err
		}

		if !hasMoreResults {
			return nil
		}
	}
	return nil
}

// newRepositoryOrgInfoResult calls the repositoryOrgInfoQuery and incorporates the respective variables
func newRepositoryOrgInfoResult(ctx context.Context, ghc ghGraphQLClient, uri string) (repository.Organization, error) {
	formattedURI, err := createNewGitHubV4URI(uri)
	if err != nil {
		return repository.Organization{}, err
	}

	repoInfoVariables := map[string]interface{}{
		"url": githubv4.URI(*formattedURI),
	}

	queryResult := new(repositoryOrgInfoQuery)
	if err := ghc.Query(ctx, queryResult, repoInfoVariables); err != nil {
		return repository.Organization{}, fmt.Errorf("RepositoryInfo query could not be completed. Error: %s", err)
	}
	if queryResult.Resource.Repository.Owner.Typename != organizationType {
		return repository.Organization{}, repository.NewTypeMismatchError(organizationType, queryResult.Resource.Repository.Owner.Typename)
	}
	organization := queryResult.Resource.Repository.Owner.Organization

	return repository.NewOrganization(organization.Name, organization.URL), nil
}

// newCommitInfoResult calls the commitInfoQuery and populates the respective variables
func newCommitInfoResult(ctx context.Context, ghc ghGraphQLClient, commitURL string) (repository.Commit, error) {
	formattedURI, err := createNewGitHubV4URI(commitURL)
	if err != nil {
		return repository.Commit{}, err
	}
	queryResult := new(commitInfoQuery)
	checkSuites, err := getAllCheckSuites(ctx, ghc, queryResult, githubv4.URI(*formattedURI))
	if err != nil {
		return repository.Commit{}, err
	}
	associatedPullRequests, err := getAllAssociatedPullRequests(ctx, ghc, queryResult, githubv4.URI(*formattedURI))
	if err != nil {
		return repository.Commit{}, err
	}
	return createNewCommitInfo(queryResult, checkSuites, associatedPullRequests)
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
		resourceType := v.(*commitInfoQuery).Resource.Typename
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
		resourceType := v.(*commitInfoQuery).Resource.Typename
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

// createNewCommitInfo returns a populated repository.CommitInfo object
func createNewCommitInfo(queryResult *commitInfoQuery, checkSuites []checkSuite, associatedPullRequests []pullRequest) (repository.Commit, error) {
	commit := queryResult.Resource.Commit

	statusState := commit.Status.State
	if !statusState.isValidStatusState() {
		return repository.Commit{}, newTypeMismatchError("statusState", statusState)
	}
	checks := make([]repository.Check, 0)
	pullRequests := make([]repository.PullRequest, 0)

	for _, pr := range associatedPullRequests {
		pullRequests = append(pullRequests, repository.NewPullRequest(pr.BaseRefName, pr.HeadRefName, pr.Merged))
	}
	for _, checkSuite := range checkSuites {
		if !checkSuite.Status.isValidCheckStatusState() {
			return repository.Commit{}, newTypeMismatchError("checkStatusState", checkSuite.Status)
		}
		if !checkSuite.Conclusion.isValidCheckConclusionState() {
			return repository.Commit{}, newTypeMismatchError("checkConclusionState", checkSuite.Conclusion)
		}
		check := repository.NewCheck(checkSuite.App.Name, checkSuite.App.URL, string(checkSuite.Status), string(checkSuite.Conclusion))
		checks = append(checks, check)
	}

	isSigned := commit.Signature.IsValid

	return repository.NewCommit(commit.URL, checks, string(statusState), isSigned, pullRequests), nil
}

// newDefaultBranchResult calls the defaultBranchQuery and populates the results with the respective variables
func newDefaultBranchResult(ctx context.Context, ghc ghGraphQLClient, repoURL string) (repository.Branch, error) {
	formattedURI, err := createNewGitHubV4URI(repoURL)
	if err != nil {
		return repository.Branch{}, err
	}
	queryResult := new(defaultBranchQuery)
	allDefaultBranchCommits := make([]commit, 0)
	defaultBranchInfoVariables := map[string]interface{}{
		"url":                       githubv4.URI(*formattedURI),
		"defaultBranchCommitCursor": (*githubv4.String)(nil),
	}

	err = paginationQuery(ctx, ghc, queryResult, defaultBranchInfoVariables, queryPageLimit, func(v interface{}) (bool, error) {
		dbq, ok := v.(*defaultBranchQuery)
		if !ok {
			return false, newTypeMismatchError("defaultBranchQuery", dbq)
		}
		resourceType := v.(*defaultBranchQuery).Resource.Typename
		if resourceType != repositoryType {
			return false, repository.NewTypeMismatchError(repositoryType, resourceType)
		}
		repo := dbq.Resource.Repository

		allDefaultBranchCommits = append(allDefaultBranchCommits, repo.DefaultBranchRef.Target.Commit.History.Nodes...)
		hasMoreResults := repo.DefaultBranchRef.Target.Commit.History.PageInfo.HasNextPage
		defaultBranchInfoVariables["defaultBranchCommitCursor"] = githubv4.NewString(repo.DefaultBranchRef.Target.Commit.History.PageInfo.EndCursor)
		return hasMoreResults, nil
	})
	if err != nil {
		return repository.Branch{}, err
	}

	defaultBranchName := queryResult.Resource.Repository.DefaultBranchRef.Name
	commits := make([]repository.CommitRef, 0)
	for _, commit := range allDefaultBranchCommits {
		repoCommit := repository.NewCommitRef(commit.URL)
		commits = append(commits, repoCommit)
	}
	return repository.NewBranch(defaultBranchName, commits), nil
}

// newBranchResult calls the branchQuery and populates the results with the respective variables
func newBranchResult(ctx context.Context, ghc ghGraphQLClient, repoURL string, branchName string) (repository.Branch, error) {
	formattedURI, err := createNewGitHubV4URI(repoURL)
	if err != nil {
		return repository.Branch{}, err
	}
	queryResult := new(branchQuery)
	allBranchCommits := make([]commit, 0)
	branchInfoVariables := map[string]interface{}{
		"url":                githubv4.URI(*formattedURI),
		"branchCommitCursor": (*githubv4.String)(nil),
		"branch_name":        (githubv4.String)(branchName),
	}

	err = paginationQuery(ctx, ghc, queryResult, branchInfoVariables, queryPageLimit, func(v interface{}) (bool, error) {
		dbq, ok := v.(*branchQuery)
		if !ok {
			return false, newTypeMismatchError("branchQuery", dbq)
		}
		resourceType := v.(*branchQuery).Resource.Typename
		if resourceType != repositoryType {
			return false, repository.NewTypeMismatchError(repositoryType, resourceType)
		}
		repo := dbq.Resource.Repository

		allBranchCommits = append(allBranchCommits, repo.Ref.Target.Commit.History.Nodes...)
		hasMoreResults := repo.Ref.Target.Commit.History.PageInfo.HasNextPage
		branchInfoVariables["branchCommitCursor"] = githubv4.NewString(repo.Ref.Target.Commit.History.PageInfo.EndCursor)
		return hasMoreResults, nil
	})
	if err != nil {
		return repository.Branch{}, err
	}

	branchNameResult := queryResult.Resource.Repository.Ref.Name
	commits := make([]repository.CommitRef, 0)
	for _, commit := range allBranchCommits {
		repoCommit := repository.NewCommitRef(commit.URL)
		commits = append(commits, repoCommit)
	}
	return repository.NewBranch(branchNameResult, commits), nil
}

// checkSuite is a collection of the check runs created by a CI/CD App
type checkSuite struct {
	App        app
	Status     checkStatusState
	Conclusion checkConclusionState
}

// checkApp contains the relevant information associated with a CI/CD app
type app struct {
	Name string
	URL  string
}

// pullRequest contains the relevant information associated with a pull request
type pullRequest struct {
	Merged      bool
	BaseRefName string
	HeadRefName string
}

// commit contains information pertaining to a commit
type commit struct {
	URL string
}
