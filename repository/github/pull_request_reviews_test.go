package github

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllReviews(t *testing.T) {
	testCases := []struct {
		testName                 string
		pullRequestURL           string
		input                    *pullRequestReviewsQuery
		mask                     []string
		queryPopulationVariables map[string]interface{}
		expected                 []review
	}{
		{
			testName:       "Testing zero associated reviews",
			pullRequestURL: "https://github.com/grafeas/voucher/pull/64",
			input: func() *pullRequestReviewsQuery {
				res := new(pullRequestReviewsQuery)
				res.Resource.Typename = "PullRequest"
				res.Resource.PullRequest.Reviews.Nodes = []review{}
				res.Resource.PullRequest.Reviews.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.PullRequest.Reviews.Nodes",
			},
			expected: []review{},
		},
		{
			testName:       "Testing has associated reviews",
			pullRequestURL: "https://github.com/grafeas/voucher/pull/23",
			input: func() *pullRequestReviewsQuery {
				res := new(pullRequestReviewsQuery)
				res.Resource.Typename = "PullRequest"
				res.Resource.PullRequest.Reviews.Nodes = []review{
					{State: "Accepted"},
					{State: "PENDING"},
				}
				res.Resource.PullRequest.Reviews.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.PullRequest.Reviews.Nodes",
			},
			expected: []review{
				{State: "Accepted"},
				{State: "PENDING"},
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.input, test.mask)
			require.Equal(t, pullRequestType, test.input.Resource.Typename)

			res, err := getAllReviews(context.Background(), c, test.pullRequestURL)
			assert.NoError(t, err, "Getting all associated reviews failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}
