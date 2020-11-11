package config

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/grafeas/voucher/repository"
	"github.com/grafeas/voucher/repository/github"
)

func TestValidRepo(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test_repo.ejson")
	viper.Set("ejson.dir", "../../testdata/key")
	viper.Set("repositories", []interface{}{map[string]interface{}{"alias": "shopify", "org-url": "github.com/Shopify"}})
	repoURL := "https://github.com/Shopify/my-app"
	secrets, err := ReadSecrets()
	assert.NoError(t, err)

	client, err := NewRepositoryClient(context.Background(), secrets.RepositoryAuthentication, repoURL)
	assert.NoError(t, err)

	assert.True(t, github.IsGithubRepoClient(client), "received client is not a github client for ", repoURL)
}

func TestInvalidRepo(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test_repo.ejson")
	viper.Set("ejson.dir", "../../testdata/key")
	viper.Set("repositories", []interface{}{map[string]interface{}{"alias": "shopify", "org-url": "github.com/Shopify"}})
	repoURL := "https://gitlab.com/TestOrg/my-app"
	secrets, err := ReadSecrets()
	assert.NoError(t, err)

	client, _ := NewRepositoryClient(context.Background(), secrets.RepositoryAuthentication, repoURL)

	assert.False(t, github.IsGithubRepoClient(client), "received github client with invalid org for ", repoURL)
}

func TestGetOrgAlias(t *testing.T) {
	orgs := map[string]repository.Organization{
		"apple":  {Alias: "apple", VCS: "github.com", Name: "my-org"},
		"banana": {Alias: "banana", VCS: "github.com", Name: ""},
	}

	cases := []struct {
		str                string
		expectedAlias      string
		expectedFoundMatch bool
	}{
		{str: "github.com/my-org/my-repo", expectedAlias: "apple", expectedFoundMatch: true},
		{str: "github.com/my-org", expectedAlias: "apple", expectedFoundMatch: true},
		{str: "github.com/other-org", expectedAlias: "banana", expectedFoundMatch: true},
		{str: "github.com", expectedAlias: "banana", expectedFoundMatch: true},
		{str: "gitea.com/hello", expectedAlias: "", expectedFoundMatch: false},
	}

	for _, testCase := range cases {
		t.Run(testCase.str, func(t *testing.T) {
			matchingKey, foundMatch := getOrgAlias(orgs, *repository.NewOrganization("", testCase.str))
			assert.Equal(t, testCase.expectedFoundMatch, foundMatch)
			assert.Equal(t, testCase.expectedAlias, matchingKey)
		})
	}
}
