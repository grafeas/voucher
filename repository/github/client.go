package github

import (
	"context"
	"fmt"

	"github.com/Shopify/voucher/repository"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// ghGraphQLClient represents a GraphQL client to interact with the GitHub API
type ghGraphQLClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// client represents the GitHub implementation of repository.Client
type client struct {
	ghClient ghGraphQLClient
}

// NewClient creates a new GitHub client
func NewClient(ctx context.Context, auth *repository.Auth) (repository.Client, error) {
	if auth == nil {
		return nil, fmt.Errorf("Must provide authentication")
	}
	if auth.Type() != repository.TokenAuthType {
		return nil, fmt.Errorf("Unsupported auth type: %s", auth.Type())
	}

	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: auth.Token},
	)
	httpClient := oauth2.NewClient(ctx, sts)
	rtw := newRoundTripperWrapper(httpClient.Transport)
	httpClient.Transport = rtw
	return &client{
		ghClient: githubv4.NewClient(httpClient),
	}, nil
}

// GetOrganization retrieves the necessary GitHub organizational information used in Voucher's checks
func (ghc *client) GetOrganization(ctx context.Context, uri string) (repository.Organization, error) {
	repoInfo, err := newRepositoryOrgInfoResult(ctx, ghc.ghClient, uri)
	if err != nil {
		return repository.Organization{}, err
	}

	if repoInfo.Resource.Typename != commitType {
		return repository.Organization{}, repository.NewTypeMismatchError(commitType, repoInfo.Resource.Typename)
	}
	if repoInfo.Resource.Commit.Repository.Owner.Typename != organizationType {
		return repository.Organization{}, repository.NewTypeMismatchError(organizationType, repoInfo.Resource.Commit.Repository.Owner.Typename)
	}

	organization := repoInfo.Resource.Commit.Repository.Owner.Organization

	return repository.CreateNewOrganization(organization.ID, organization.Name, organization.URL), nil
}

// GetCommitInfo retrieves the necessary GitHub commit information used in Voucher's checks
func (ghc *client) GetCommitInfo(ctx context.Context, commitURI string) (repository.CommitInfo, error) {
	commitInfo, err := newCommitInfoResult(ctx, ghc.ghClient, commitURI)
	if err != nil {
		return repository.CommitInfo{}, fmt.Errorf("GetCommitInfo query could not be completed. Error: %s", err)
	}
	return commitInfo, nil
}

// GetDefaultBranch retrieves the necessary GitHub default branch information used in Voucher's checks
func (ghc *client) GetDefaultBranch(ctx context.Context, commitURI string) (repository.DefaultBranch, error) {
	defaultBranchResult, err := newDefaultBranchResult(ctx, ghc.ghClient, commitURI)
	if err != nil {
		return repository.DefaultBranch{}, fmt.Errorf("GetDefaultBranch query could not be completed. Error: %s", err)
	}
	return defaultBranchResult, nil
}
