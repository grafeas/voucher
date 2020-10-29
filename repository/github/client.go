package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/grafeas/voucher/repository"
)

// errCreatingRepositoryMetadata is the error returned when we fail to create
// repository metadata.
var errCreatingRepositoryMetadata = errors.New("failed to create repository metadata")

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
		return nil, fmt.Errorf("must provide authentication")
	}

	var httpClient *http.Client
	if auth.Type() == repository.TokenAuthType {
		sts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: auth.Token},
		)
		httpClient = oauth2.NewClient(ctx, sts)
		rtw := newRoundTripperWrapper(httpClient.Transport)
		httpClient.Transport = rtw
	} else if auth.Type() == repository.GithubInstallType {
		appID, err := strconv.Atoi(auth.AppID)
		if err != nil {
			return nil, fmt.Errorf("invalid application ID: %v", err)
		}

		installID, err := strconv.Atoi(auth.InstallationID)
		if err != nil {
			return nil, fmt.Errorf("invalid installation ID: %v", err)
		}

		appsTransport, err := ghinstallation.New(http.DefaultTransport, int64(appID), int64(installID), []byte(auth.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("error configuring Github App transport: %v", err)
		}
		httpClient = &http.Client{}
		httpClient.Transport = newRoundTripperWrapper(appsTransport)
	} else {
		return nil, fmt.Errorf("unsupported auth type: %s", auth.Type())
	}

	return &client{
		ghClient: githubv4.NewClient(httpClient),
	}, nil
}

// GetOrganization retrieves the necessary GitHub organizational information used in Voucher's checks
func (ghc *client) GetOrganization(ctx context.Context, details repository.BuildDetail) (repository.Organization, error) {
	repo := repository.NewRepositoryMetadata(details.RepositoryURL)
	if nil == repo {
		return repository.Organization{}, fmt.Errorf("error parsing repository url %s", details.RepositoryURL)
	}

	organization, err := newRepositoryOrgInfoResult(ctx, ghc.ghClient, repo.String())
	if err != nil {
		return repository.Organization{}, err
	}

	return organization, nil
}

// GetCommit retrieves the necessary GitHub commit information used in Voucher's checks
func (ghc *client) GetCommit(ctx context.Context, details repository.BuildDetail) (repository.Commit, error) {
	commitURI, err := GetCommitURL(&details)
	if err != nil {
		return repository.Commit{}, fmt.Errorf("error creating a commit url. Error: %s", err)
	}

	commit, err := newCommitInfoResult(ctx, ghc.ghClient, commitURI)
	if err != nil {
		return repository.Commit{}, fmt.Errorf("GetCommitInfo query could not be completed. Error: %s", err)
	}
	return commit, nil
}

// GetDefaultBranch retrieves the necessary GitHub default branch information used in Voucher's checks
func (ghc *client) GetDefaultBranch(ctx context.Context, details repository.BuildDetail) (repository.Branch, error) {
	repo := repository.NewRepositoryMetadata(details.RepositoryURL)
	if nil == repo {
		return repository.Branch{}, errCreatingRepositoryMetadata
	}

	defaultBranchResult, err := newDefaultBranchResult(ctx, ghc.ghClient, repo.String())
	if err != nil {
		return repository.Branch{}, fmt.Errorf("GetDefaultBranch query could not be completed. Error: %s", err)
	}
	return defaultBranchResult, nil
}

// GetBranch retrieves the necessary GitHub branch information used in Voucher's checks given the name of the branch
func (ghc *client) GetBranch(ctx context.Context, details repository.BuildDetail, name string) (repository.Branch, error) {
	repo := repository.NewRepositoryMetadata(details.RepositoryURL)
	if nil == repo {
		return repository.Branch{}, errCreatingRepositoryMetadata
	}

	branchResult, err := newBranchResult(ctx, ghc.ghClient, repo.String(), name)
	if err != nil {
		return repository.Branch{}, fmt.Errorf("GetBranch query could not be completed. Error: %s", err)
	}
	return branchResult, nil
}

func IsGithubRepoClient(repositoryClient repository.Client) bool {
	_, ok := repositoryClient.(*client)
	return ok
}
