package config

import (
	"context"
	"fmt"

	"github.com/grafeas/voucher/repository"
	"github.com/grafeas/voucher/repository/github"
)

// NewRepositoryClient creates a new repository.Client for the given repository URL. The URL may be in any known
// format including, but not limited to, urls starting with 'http://', 'https://', 'git@', etc.
func NewRepositoryClient(ctx context.Context, keyring repository.KeyRing, repoURL string) (repository.Client, error) {
	org := repository.NewOrganization("", repoURL)
	if nil == org {
		return nil, fmt.Errorf("error parsing url %s", repoURL)
	}

	token, err := getTokenForOrg(keyring, *org)
	if nil != err {
		return nil, err
	}

	switch org.VCS {
	case "github.com":
		return github.NewClient(context.Background(), token)
	}

	return nil, fmt.Errorf("unknown repository %s", repoURL)
}

func getTokenForOrg(keyring repository.KeyRing, org repository.Organization) (*repository.Auth, error) {
	orgs := GetOrganizationsFromConfig()
	if alias, ok := getOrgAlias(orgs, org); ok {
		token := keyring[alias]
		return &token, nil
	}

	return nil, fmt.Errorf("failed to get token for %s", org.Alias)
}

func getOrgAlias(orgs map[string]repository.Organization, repoOrg repository.Organization) (matchingKey string, foundMatch bool) {
	var matchLength int
	var longestMatch string

	for alias, org := range orgs {
		if !isMatch(org, repoOrg) {
			continue
		}

		// catch all
		if 1 > matchLength {
			longestMatch = alias
		}

		if org.VCS == repoOrg.VCS && 2 > matchLength {
			matchLength = 1
			longestMatch = alias
		}

		if org.Name == repoOrg.Name {
			matchLength = 2
			longestMatch = alias
		}
	}

	return longestMatch, "" != longestMatch
}

func isMatch(org, repoOrg repository.Organization) bool {
	if org.VCS == "" && org.Name == "" {
		return true
	}

	if org.VCS != repoOrg.VCS {
		return false
	}

	return org.Name == "" || org.Name == repoOrg.Name
}
