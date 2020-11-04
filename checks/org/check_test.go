package org

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/repository"
	r "github.com/grafeas/voucher/repository"
)

func TestOrgCheck(t *testing.T) {
	c := context.Background()

	i, err := voucher.NewImageData("gcr.io/voucher-test-project/apps/staging/voucher-internal@sha256:73d506a23331fce5cb6f49bfb4c27450d2ef4878efce89f03a46b27372a88430")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)
	details := r.BuildDetail{RepositoryURL: "https://github.com/Shopify/app", Commit: "efgh6543"}
	organization := r.Organization{Name: "Shopify", VCS: "github.com"}

	repoClient := new(r.MockClient)
	repoClient.On("GetOrganization", mock.Anything, details).Return(organization, nil)

	metadataClient := new(voucher.MockMetadataClient)
	metadataClient.On("GetBuildDetail", mock.Anything, i).Return(details, nil)

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

	i, err := voucher.NewImageData("gcr.io/voucher-test-project/apps/staging/voucher-internal@sha256:73d506a23331fce5cb6f49bfb4c27450d2ef4878efce89f03a46b27372a88430")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)
	details := r.BuildDetail{RepositoryURL: "git@github.com/TestOrg/TestRepo.git", Commit: "cdef0987"}
	organization := r.Organization{Name: "Shopify", VCS: "github.com"}

	repoClient := new(r.MockClient)
	repoClient.On("GetOrganization", mock.Anything, details).Return(repository.Organization{}, nil)

	metadataClient := new(voucher.MockMetadataClient)
	metadataClient.On("GetBuildDetail", mock.Anything, i).Return(details, nil)

	orgCheck := new(check)
	orgCheck.org = organization
	orgCheck.SetRepositoryClient(repoClient)
	orgCheck.SetMetadataClient(metadataClient)

	status, err := orgCheck.Check(c, i)

	assert.NoErrorf(t, err, "check failed with error: %s", err)
	assert.False(t, status, "check passed when it should have failed")
}
