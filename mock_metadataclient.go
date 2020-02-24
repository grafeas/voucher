package voucher

import (
	"context"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/mock"

	"github.com/grafeas/voucher/repository"
)

type MockMetadataClient struct {
	mock.Mock
}

func (m *MockMetadataClient) CanAttest() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockMetadataClient) NewPayloadBody(ref reference.Canonical) (string, error) {
	args := m.Called(ref)
	return args.String(0), args.Error(1)
}

func (m *MockMetadataClient) GetVulnerabilities(ctx context.Context, ref reference.Canonical) ([]Vulnerability, error) {
	args := m.Called(ctx, ref)
	return args.Get(0).([]Vulnerability), args.Error(1)
}

func (m *MockMetadataClient) GetBuildDetail(ctx context.Context, ref reference.Canonical) (repository.BuildDetail, error) {
	args := m.Called(ctx, ref)
	return args.Get(0).(repository.BuildDetail), args.Error(1)
}

func (m *MockMetadataClient) AddAttestationToImage(ctx context.Context, ref reference.Canonical, attestation Attestation) (SignedAttestation, error) {
	args := m.Called(ctx, ref, attestation)
	return args.Get(0).(SignedAttestation), args.Error(1)
}

func (m *MockMetadataClient) GetAttestations(ctx context.Context, ref reference.Canonical) ([]SignedAttestation, error) {
	args := m.Called(ctx, ref)
	return args.Get(0).([]SignedAttestation), args.Error(1)
}

func (m *MockMetadataClient) Close() {
	m.Called()
}
