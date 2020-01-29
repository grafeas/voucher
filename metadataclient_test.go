package voucher

import (
	"context"
	"errors"

	"github.com/Shopify/voucher/repository"
	"github.com/docker/distribution/reference"
)

type testMetadataClient struct {
	canAttest bool
	keyring   *KeyRing
}

func (t *testMetadataClient) Close() {
}

func (t *testMetadataClient) CanAttest() bool {
	return t.canAttest
}

func (t *testMetadataClient) NewPayloadBody(i ImageData) (string, error) {
	if t.canAttest {
		return i.String(), nil
	}
	return "", errors.New("cannot create payload body")
}

func (t *testMetadataClient) GetVulnerabilities(ctx context.Context, i ImageData) ([]Vulnerability, error) {
	return []Vulnerability{}, nil
}

func (t *testMetadataClient) AddAttestationToImage(ctx context.Context, i ImageData, payload AttestationPayload) (interface{}, error) {
	_, _, err := payload.Sign(t.keyring)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

func (t *testMetadataClient) GetBuildDetails(ctx context.Context, reference reference.Canonical) ([]repository.BuildDetail, error) {
	return []repository.BuildDetail{}, nil
}

func newTestMetadataClient(keyring *KeyRing, canAttest bool) *testMetadataClient {
	client := new(testMetadataClient)
	client.keyring = keyring
	client.canAttest = canAttest
	return client
}
