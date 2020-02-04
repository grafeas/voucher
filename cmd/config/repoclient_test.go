package config

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/Shopify/voucher/repository"
	"github.com/Shopify/voucher/repository/github"
)

func TestValidRepo(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test_repo.ejson")
	repoURL := "https://github.com/Shopify/my-app"
	client, err := NewRepositoryClient(context.Background(), repoURL)
	assert.NoError(t, err)

	assert.True(t, github.IsGithubRepoClient(client), "received client is not a github client for ", repoURL)
}

func TestInvalidRepo(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test_repo.ejson")
	repoURL := "https://gitlab.com/TestOrg/my-app"
	client, _ := NewRepositoryClient(context.Background(), repoURL)

	assert.False(t, github.IsGithubRepoClient(client), "received github client with invalid org for ", repoURL)
}

func TestFindLongestKeyMatch(t *testing.T) {
	keyRing := repository.KeyRing{
		"github.com":        repository.Auth{},
		"github.com/my-org": repository.Auth{},
	}

	cases := []struct {
		str         string
		matchingKey string
		foundMatch  bool
	}{
		{str: "github.com/my-org/my-repo", matchingKey: "github.com/my-org", foundMatch: true},
		{str: "github.com/my-org", matchingKey: "github.com/my-org", foundMatch: true},
		{str: "github.com/other-org", matchingKey: "github.com", foundMatch: true},
		{str: "gitea.com/hello", matchingKey: "", foundMatch: false},
	}

	for _, testCase := range cases {
		t.Run(testCase.str, func(t *testing.T) {
			matchingKey, foundMatch := findLongestKeyMatch(keyRing, testCase.str)
			assert.Equal(t, testCase.foundMatch, foundMatch)
			assert.Equal(t, testCase.matchingKey, matchingKey)
		})
	}
}
