package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/grafeas/voucher/repository"
	"github.com/shurcooL/githubv4"
)

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

	org := repository.NewOrganization(organization.Name, organization.URL)
	if org == nil {
		return repository.Organization{}, errors.New("error parsing url" + organization.URL)
	}

	return *org, nil
}
