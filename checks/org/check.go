package org

import (
	"context"
	"errors"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/repository"
)

// ErrNoBuildData is an error returned if we can't pull any BuildData from
// Grafeas for an image.
var ErrNoBuildData = errors.New("no build metadata associated with this image")

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

// Check runs the org check
func (o *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	items, err := o.metadataClient.GetBuildDetails(ctx, i)
	if err != nil {
		if voucher.IsNoMetadataError(err) {
			return false, ErrNoBuildData
		}
		return false, err
	}

	for _, buildDetail := range items {
		org, err := o.repositoryClient.GetOrganization(ctx, buildDetail)
		if err != nil {
			return false, err
		}
		if org != o.org {
			return false, nil
		}
	}

	return true, nil
}

func NewOrganizationCheckFactory(organization repository.Organization, client repository.Client) voucher.CheckFactory {
	// Return a voucher.CheckFactory that always creates the desired OrganizationCheck
	return func() voucher.Check {
		return &check{
			org:              organization,
			repositoryClient: client,
		}
	}
}
