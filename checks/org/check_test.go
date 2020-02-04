package org

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Shopify/voucher"
	r "github.com/Shopify/voucher/repository"
	vtesting "github.com/Shopify/voucher/testing"
)

func TestOrgCheck(t *testing.T) {
	c := context.Background()

	repoClient, _ := vtesting.NewClient()
	commits := map[string]r.Commit{
		"cdef0987": {URL: "https://github.com/Shopify/app/commit/cdef0987", Checks: nil, Status: "good", IsSigned: true,
			AssociatedPullRequests: nil},
		"efgh6543": {URL: "https://github.com/Shopify/app/commit/efgh6543", Checks: nil, Status: "great", IsSigned: false,
			AssociatedPullRequests: []r.PullRequest{{BaseBranchName: "base", HeadBranchName: "master", IsMerged: false},
				{BaseBranchName: "branch-base", HeadBranchName: "test", IsMerged: true}}},
	}
	branches := map[string]r.Branch{
		"test":   {Name: "test", CommitRefs: []r.CommitRef{{URL: "https://github.com/Shopify/app/commit/efgh6543"}}},
		"master": {Name: "master", CommitRefs: []r.CommitRef{{URL: "https://github.com/Shopify/app/commit/cdef0987"}}},
	}
	repoClient.AddRepository(r.NewOrganization("Shopify", "https://github.com/Shopify"), "Repo_Shopify", commits, branches)

	metadataClient := new(vtesting.MetadataClient)
	i, err := voucher.NewImageData("gcr.io/voucher-test-project/apps/staging/voucher-internal@sha256:73d506a23331fce5cb6f49bfb4c27450d2ef4878efce89f03a46b27372a88430")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)
	details := r.BuildDetail{RepositoryURL: "https://github.com/Shopify/app", Commit: "efgh6543"}
	metadataClient.SetBuildDetail(i, details)

	organization := r.Organization{Name: "Shopify", URL: "https://github.com/Shopify"}

	orgCheck := new(check)
	orgCheck.org = organization
	orgCheck.SetRepositoryClient(repoClient)
	orgCheck.SetMetadataClient(metadataClient)

	status, err := orgCheck.Check(c, i)

	assert.NoErrorf(t, err, "check failed with error: %s", err)
	assert.True(t, status, "check failed when it should have passed")
}

func TestOrgCheckWithInvalidRepo(t *testing.T) {
	c := context.Background()

	repoClient, _ := vtesting.NewClient()
	commits := map[string]r.Commit{
		"cdef0987": {URL: "https://github.com/TestOrg/app/commit/cdef0987", Checks: nil, Status: "good", IsSigned: true,
			AssociatedPullRequests: nil},
	}
	branches := map[string]r.Branch{
		"master": {Name: "master", CommitRefs: []r.CommitRef{{URL: "https://github.com/TestOrg/app/commit/cdef0987"}}},
	}
	repoClient.AddRepository(r.NewOrganization("TestOrg", "https://github.com/TestOrg"), "TestRepo", commits, branches)

	metadataClient := new(vtesting.MetadataClient)
	i, err := voucher.NewImageData("gcr.io/voucher-test-project/apps/staging/voucher-internal@sha256:73d506a23331fce5cb6f49bfb4c27450d2ef4878efce89f03a46b27372a88430")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)
	details := r.BuildDetail{RepositoryURL: "git@github.com/TestOrg/TestRepo.git", Commit: "cdef0987"}
	metadataClient.SetBuildDetail(i, details)

	organization := r.Organization{Name: "Shopify", URL: "https://github.com/Shopify"}

	orgCheck := new(check)
	orgCheck.org = organization
	orgCheck.SetRepositoryClient(repoClient)
	orgCheck.SetMetadataClient(metadataClient)

	status, err := orgCheck.Check(c, i)

	assert.NoErrorf(t, err, "check failed with error: %s", err)
	assert.False(t, status, "check passed when it should have failed")
}
