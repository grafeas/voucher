package org

import (
	"context"
	"errors"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/repository"
)

// ErrNoBuildData is an error returned if we can't pull any BuildData from
// Grafeas for an image.
var ErrNoBuildData = errors.New("no build metadata associated with this image")

// ErrNoRepositoryClient is an error returned if we can't connect to the source code repository for an image.
var ErrNoRepositoryClient = errors.New("no repository client configured for check")

// check holds the required data for the check
type check struct {
	metadataClient   voucher.MetadataClient
	repositoryClient repository.Client
	org              repository.Organization
}

// SetMetadataClient sets the MetadataClient for this Check.
func (o *check) SetMetadataClient(metadataClient voucher.MetadataClient) {
	o.metadataClient = metadataClient
}

// SetRepositoryClient sets the repository client for this Check.
func (o *check) SetRepositoryClient(repositoryClient repository.Client) {
	o.repositoryClient = repositoryClient
}

// Check runs the org check
func (o *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	buildDetail, err := o.metadataClient.GetBuildDetail(ctx, i)
	if err != nil {
		if voucher.IsNoMetadataError(err) {
			return false, ErrNoBuildData
		}
		return false, err
	}

	if o.repositoryClient == nil {
		return false, ErrNoRepositoryClient
	}

	org, err := o.repositoryClient.GetOrganization(ctx, buildDetail)
	if err != nil {
		return false, err
	}
	if org.Name != o.org.Name {
		return false, nil
	}

	return true, nil
}

func NewOrganizationCheckFactory(organization repository.Organization) voucher.CheckFactory {
	// Return a voucher.CheckFactory that always creates the desired OrganizationCheck
	return func() voucher.Check {
		return &check{
			org: organization,
		}
	}
}
