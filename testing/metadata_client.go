package vtesting

import (
	"context"
	"errors"

	"github.com/docker/distribution/reference"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/attestation"
	"github.com/Shopify/voucher/repository"
)

type MetadataClient struct {
	details map[string][]repository.BuildDetail
	vulns   map[string][]voucher.Vulnerability
	keyring *voucher.KeyRing
}

//AddBuildDetail adds BuildDetails to the metadata client
func (c *MetadataClient) AddBuildDetails(reference reference.Reference, details []repository.BuildDetail) {
	if c.details == nil {
		c.details = make(map[string][]repository.BuildDetail)
	}
	refString := reference.String()
	c.details[refString] = append(c.details[refString], details...)
}

//AddBuildDetail adds Vulnerabilities to the metadata client
func (c *MetadataClient) AddVulnerabilities(reference reference.Reference, vulnerabilities []voucher.Vulnerability) {
	if c.vulns == nil {
		c.vulns = make(map[string][]voucher.Vulnerability)
	}
	refString := reference.String()
	c.vulns[refString] = append(c.vulns[refString], vulnerabilities...)
}

func (c *MetadataClient) CanAttest() bool {
	return nil != c.keyring
}

func (c *MetadataClient) NewPayloadBody(reference reference.Canonical) (string, error) {
	payload, err := attestation.NewPayload(reference).ToString()
	if err != nil {
		return "", err
	}
	return payload, err
}

func (c *MetadataClient) GetVulnerabilities(ctx context.Context, reference reference.Canonical) ([]voucher.Vulnerability, error) {
	refString := reference.String()
	return c.vulns[refString], nil
}

func (c *MetadataClient) GetBuildDetails(ctx context.Context, reference reference.Canonical) ([]repository.BuildDetail, error) {
	refString := reference.String()
	if len(c.details) == 0 {
		err := &voucher.NoMetadataError{
			Type: voucher.VulnerabilityType,
			Err:  errors.New("no occurrences returned for image"),
		}
		return nil, err
	}
	return c.details[refString], nil
}

func (c *MetadataClient) AddAttestationToImage(ctx context.Context, reference reference.Canonical, payload voucher.AttestationPayload) (interface{}, error) {
	return nil, nil
}

func (c *MetadataClient) Close() {}
