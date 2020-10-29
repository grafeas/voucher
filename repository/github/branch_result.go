package github

import (
	"context"

	"github.com/grafeas/voucher/repository"
	"github.com/shurcooL/githubv4"
)

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
