package github

import (
	"context"

	"github.com/grafeas/voucher/repository"
	"github.com/shurcooL/githubv4"
)

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
