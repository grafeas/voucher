package github

import (
	"testing"

	"github.com/grafeas/voucher/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetCommitURL(t *testing.T) {
	getCommitURLTests := []struct {
		expectedURL     string
		mockBuildDetail *repository.BuildDetail
	}{
		{
			expectedURL: "https://github.com/grafeas/voucher/commit/sl2o3vo2wojweoie",
			mockBuildDetail: &repository.BuildDetail{
				RepositoryURL: "git@github.com/grafeas/voucher.git",
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
