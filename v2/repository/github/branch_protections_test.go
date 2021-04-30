package github

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBranchProtections(t *testing.T) {
	testCases := []struct {
		testName                 string
		repoURL                  string
		input                    *branchProtectionsQuery
		mask                     []string
		queryPopulationVariables map[string]interface{}
		expected                 map[string]branchProtection
	}{
		{
			testName: "Testing zero associated branch protections",
			repoURL:  "https://github.com/grafeas/voucher",
			input: func() *branchProtectionsQuery {
				res := new(branchProtectionsQuery)
				res.Resource.Typename = "Repository"
				res.Resource.Repository.BranchProtectionRules.Nodes = []branchProtectionRule{}
				res.Resource.Repository.BranchProtectionRules.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Repository.BranchProtectionRules.Nodes",
			},
			expected: make(map[string]branchProtection),
		},
		{
			testName: "Testing zero associated matchingRefs for branch protections",
			repoURL:  "https://github.com/grafeas/voucher",
			input: func() *branchProtectionsQuery {
				res := new(branchProtectionsQuery)
				res.Resource.Typename = "Repository"
				protectionRule := branchProtectionRule{}
				protectionRule.MatchingRefs.Nodes = []matchingRef{}
				protectionRule.MatchingRefs.PageInfo.HasNextPage = false
				res.Resource.Repository.BranchProtectionRules.Nodes = []branchProtectionRule{protectionRule}
				res.Resource.Repository.BranchProtectionRules.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Repository.BranchProtectionRules.Nodes",
			},
			expected: make(map[string]branchProtection),
		},
		{
			testName: "Testing has associated branch protections",
			repoURL:  "https://github.com/grafeas/voucher",
			input: func() *branchProtectionsQuery {
				res := new(branchProtectionsQuery)
				res.Resource.Typename = "Repository"
				protectionRule := branchProtectionRule{
					RequiresApprovingReviews:     true,
					RequiredApprovingReviewCount: 2,
				}
				protectionRule.MatchingRefs.Nodes = []matchingRef{
					{Name: "master"},
					{Name: "staging"},
				}
				protectionRule.MatchingRefs.PageInfo.HasNextPage = false
				res.Resource.Repository.BranchProtectionRules.Nodes = []branchProtectionRule{protectionRule}
				res.Resource.Repository.BranchProtectionRules.PageInfo.HasNextPage = false
				return res
			}(),
			mask: []string{
				"Resource.Typename",
				"Resource.Repository.BranchProtectionRules.Nodes",
			},
			expected: map[string]branchProtection{
				"master": {
					RequiresApprovingReviews:     true,
					RequiredApprovingReviewCount: 2,
				},
				"staging": {
					RequiresApprovingReviews:     true,
					RequiredApprovingReviewCount: 2,
				},
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			c := new(mockGitHubGraphQLClient)
			c.HandlerFunc = createHandler(test.input, test.mask)
			require.Equal(t, repositoryType, test.input.Resource.Typename)

			res, err := getBranchProtections(context.Background(), c, test.repoURL)
			assert.NoError(t, err, "Getting all associated branch protections failed")
			assert.EqualValues(t, test.expected, res)
		})
	}
}
