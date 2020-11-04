package github

import (
	"context"
	"testing"

	"github.com/grafeas/voucher/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepositoryOrgInfoResult(t *testing.T) {
	newRepoOrgTests := []struct {
		testName    string
		uri         string
		input       *repositoryOrgInfoQuery
		mask        []string
		expected    repository.Organization
		shouldError bool
	}{
		{
			testName: "Testing happy path",
			uri:      "https://github.com/grafeas/voucher",
			input: func() *repositoryOrgInfoQuery {
				res := new(repositoryOrgInfoQuery)
				res.Resource.Repository.Owner.Typename = "Organization"
				res.Resource.Repository.Owner.Organization.Name = "Shopify"
				res.Resource.Repository.Owner.Organization.URL = "https://github.com/Shopify"
				return res
			}(),
			mask: []string{"Resource.Repository.Owner.Typename", "Resource.Repository.Owner.Organization"},
			expected: repository.Organization{
				Alias: "Shopify",
				VCS:   "github.com",
				Name:  "Shopify",
			},
			shouldError: false,
		},
		{
			testName:    "Testing with bad URL",
			uri:         "hello@%a&%(.com",
			input:       new(repositoryOrgInfoQuery),
			mask:        []string{},
			expected:    repository.Organization{},
			shouldError: true,
		},
	}

	for _, test := range newRepoOrgTests {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.input, test.mask)
			res, err := newRepositoryOrgInfoResult(context.Background(), c, test.uri)
			if test.shouldError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, test.expected, res)
			assert.Equal(t, organizationType, test.input.Resource.Repository.Owner.Typename)
		})
	}
}
