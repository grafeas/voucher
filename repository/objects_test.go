package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrganization(t *testing.T) {
	cases := []struct {
		alias    string
		url      string
		expected Organization
	}{
		{
			alias: "MyOrg",
			url:   "https://github.com/organization",
			expected: Organization{
				Alias: "MyOrg",
				Name:  "organization",
				VCS:   "github.com",
			},
		},
		{
			alias: "MyOrg2",
			url:   "https://github.com",
			expected: Organization{
				Alias: "MyOrg2",
				VCS:   "github.com",
			},
		},
		{
			alias: "MyOrg3",
			url:   "gitlab.com/Org/repo",
			expected: Organization{
				Alias: "MyOrg3",
				VCS:   "gitlab.com",
				Name:  "Org",
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.alias, func(t *testing.T) {
			org := NewOrganization(testCase.alias, testCase.url)
			require.NotNil(t, org)
			assert.Equal(t, testCase.expected.Alias, org.Alias)
			assert.Equal(t, testCase.expected.Name, org.Name)
			assert.Equal(t, testCase.expected.VCS, org.VCS)
		})
	}
}

func TestNewRepositoryMetadata(t *testing.T) {
	cases := []struct {
		url      string
		expected Metadata
	}{
		{
			url: "https://github.com/my-org/my-repo",
			expected: Metadata{
				Name:         "my-repo",
				Organization: "my-org",
				VCS:          "github.com",
			},
		},
		{
			url: "https://github.com/my-org",
			expected: Metadata{
				Organization: "my-org",
				VCS:          "github.com",
			},
		},
		{
			url: "gitlab.com/",
			expected: Metadata{
				VCS: "gitlab.com",
			},
		},
		{
			url: "git@github.com/my-org/my-repo.git",
			expected: Metadata{
				Name:         "my-repo",
				Organization: "my-org",
				VCS:          "github.com",
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.url, func(t *testing.T) {
			metadata := NewRepositoryMetadata(testCase.url)
			assert.Equal(t, testCase.expected.Name, metadata.Name)
			assert.Equal(t, testCase.expected.Organization, metadata.Organization)
			assert.Equal(t, testCase.expected.VCS, metadata.VCS)
		})
	}
}
