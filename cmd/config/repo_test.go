package config

import (
	"net/url"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	r "github.com/Shopify/voucher/repository"
	"github.com/Shopify/voucher/repository/github"
)

func TestValidOrg(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test_repo.ejson")
	organization := r.Organization{Name: "Shopify", URL: "https://github.com/Shopify"}
	client := newRepositoryClient(organization)

	assert.True(t, github.IsGithubRepoClient(client), "received client is not a github client: ", client)
}

func TestInvalidOrg(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test_repo.ejson")
	organization := r.Organization{Name: "TestOrg", URL: "https://gitlab.com/TestOrg"}
	client := newRepositoryClient(organization)

	assert.False(t, github.IsGithubRepoClient(client), "received github client with invalid org: "+organization.Name)
}

func TestValidFullUrl(t *testing.T) {
	orgURL, _ := url.Parse("https://github.com/Shopify")
	repoURL, _ := url.Parse("https://github.com/Shopify")
	res := isOrgURL(orgURL, repoURL)

	assert.True(t, res, "same urls were marked as different, org: "+orgURL.Path+", repo: "+repoURL.Path)
}

func TestInValidFullUrl(t *testing.T) {
	orgURL, err := url.Parse("https://gitlab.com/Shopify")
	assert.NoError(t, err)
	repoURL, err := url.Parse("https://github.com/Shopify")
	assert.NoError(t, err)
	res := isOrgURL(orgURL, repoURL)

	assert.False(t, res, "different urls were marked as same, org: "+orgURL.Path+", repo: "+repoURL.Path)
}

func TestParseValidFullUrl(t *testing.T) {
	urlStr := "https://gitlab.com/TestOrg"
	urlParsed, err := normalizeURL(urlStr)
	assert.NoError(t, err)

	assert.Equal(t, "gitlab.com", urlParsed.Hostname())
	assert.Equal(t, "/TestOrg", urlParsed.Path)
}

func TestParseValidNotFullUrl(t *testing.T) {
	urlStr := "gitlab.com/TestOrg"
	urlParsed, err := normalizeURL(urlStr)
	assert.NoError(t, err)

	assert.Equal(t, "gitlab.com", urlParsed.Hostname())
	assert.Equal(t, "/TestOrg", urlParsed.Path)
}

func TestParseInValidUrl(t *testing.T) {
	urlStr := "test%random&string"
	_, err := normalizeURL(urlStr)

	assert.Error(t, err)
}
