package github

import (
	"testing"

	"github.com/Shopify/voucher/repository"
	"github.com/stretchr/testify/assert"
)

func TestParseGithutURL(t *testing.T) {
	parseGithubTests := []struct {
		expectedResults *RepositoryMetadata
		githubURL       string
	}{
		{
			expectedResults: &RepositoryMetadata{
				Vcs:          "github.com",
				Organization: "Shopify",
				Name:         "TestRepo",
			},
			githubURL: "git@github.com/Shopify/TestRepo.git",
		},
		{
			expectedResults: &RepositoryMetadata{
				Vcs:          "github.com",
				Organization: "Shopify",
				Name:         "voucher",
			},
			githubURL: "https://github.com/Shopify/voucher",
		},
	}

	for _, test := range parseGithubTests {
		repoMeta, err := ParseGithubURL(test.githubURL)
		if assert.NoError(t, err, "Failed to parse correct Github url") {
			assert.EqualValues(t, test.expectedResults, repoMeta, "github url is not parsed properly")
		}
	}
}

func TestInvalidGithubURL(t *testing.T) {
	mockInvalidURL := "something@git.not.correct.git"
	repoMeta, err := ParseGithubURL(mockInvalidURL)
	if assert.Error(t, err, "error, an invalid github repo url is parsed") {
		assert.Nil(t, repoMeta, "error, invalid github url is parsed")
	}
}

func TestGetCommitURL(t *testing.T) {
	getCommitURLTests := []struct {
		expectedURL     string
		mockBuildDetail *repository.BuildDetail
	}{
		{
			expectedURL: "https://github.com/Shopify/voucher/commit/sl2o3vo2wojweoie",
			mockBuildDetail: &repository.BuildDetail{
				RepositoryURL: "git@github.com/Shopify/voucher.git",
				Commit:        "sl2o3vo2wojweoie",
				BuildCreator:  "someone",
				BuildURL:      "somebuild.url.io",
			},
		},
	}
	for _, test := range getCommitURLTests {
		commitURL, err := GetCommitURL(test.mockBuildDetail)
		assert.NoError(t, err, "error parsing github url")
		assert.EqualValues(t, test.expectedURL, commitURL, "commit url is not properly formatted")
	}
}

func TestGetRepositoryURL(t *testing.T) {
	getRepositoryURLTests := []struct {
		expectedURL     string
		mockBuildDetail *repository.BuildDetail
	}{
		{
			expectedURL: "https://github.com/Shopify/voucher",
			mockBuildDetail: &repository.BuildDetail{
				RepositoryURL: "git@github.com/Shopify/voucher.git",
				Commit:        "sl2o3vo2wojweoie",
				BuildCreator:  "someone",
				BuildURL:      "somebuild.url.io",
			},
		},
	}
	for _, test := range getRepositoryURLTests {
		commitURL, err := GetRepositoryURL(test.mockBuildDetail)
		assert.NoError(t, err, "error parsing github url")
		assert.EqualValues(t, test.expectedURL, commitURL, "repository url is not properly formatted")
	}
}
