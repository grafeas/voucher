package config

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/Shopify/voucher/repository"
	"github.com/Shopify/voucher/repository/github"
)

// NewRepositoryClient creates a new repository.Client for the given repository URL. The URL may be in any known
// format including, but not limited to, urls starting with 'http://', 'https://', 'git@', etc.
func NewRepositoryClient(ctx context.Context, repoURL string) (repository.Client, error) {
	parseFuncs := []func(repoURL string) (*repository.Metadata, error){
		github.ParseGithubURL,
	}

	for _, parseRepoURL := range parseFuncs {
		repo, err := parseRepoURL(repoURL)
		if nil != err {
			continue
		}

		token, err := getTokenForRepo(repo)
		if nil != err {
			return nil, err
		}

		switch repo.Vcs {
		case "github.com":
			return github.NewClient(context.Background(), token)
		}
	}

	return nil, fmt.Errorf("unable to parse repo url %s", repoURL)
}

func getTokenForRepo(repo *repository.Metadata) (*repository.Auth, error) {
	keyring, err := getRepositoryKeyRing()
	if nil != err {
		return nil, fmt.Errorf("failed to get repository key ring: %w", err)
	}

	orgURL := path.Join(repo.Vcs, repo.Organization)
	if key, ok := findLongestKeyMatch(keyring, orgURL); ok {
		token := keyring[key]
		return &token, nil
	}

	return nil, fmt.Errorf("failed to get token for %s: %w", orgURL, err)
}

func findLongestKeyMatch(keyRing repository.KeyRing, toMatch string) (matchingKey string, foundMatch bool) {
	var matchLength int
	var longestMatch string

	for key := range keyRing {
		if strings.HasPrefix(toMatch, key) {
			if len(key) > matchLength {
				matchLength = len(key)
				longestMatch = key
			}
		}
	}

	return longestMatch, matchLength > 0
}
