// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	context "context"

	cyclonedx "github.com/CycloneDX/cyclonedx-go"
	mock "github.com/stretchr/testify/mock"

	reference "github.com/docker/distribution/reference"

	testing "testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	voucher "github.com/grafeas/voucher/v2"
)

// GCRClient is an autogenerated mock type for the GCRClient type
type GCRClient struct {
	mock.Mock
}

// GetSBOM provides a mock function with given fields: ctx, imageName, tag
func (_m *GCRClient) GetSBOM(ctx context.Context, imageName string, tag string) (cyclonedx.BOM, error) {
	ret := _m.Called(ctx, imageName, tag)

	var r0 cyclonedx.BOM
	if rf, ok := ret.Get(0).(func(context.Context, string, string) cyclonedx.BOM); ok {
		r0 = rf(ctx, imageName, tag)
	} else {
		r0 = ret.Get(0).(cyclonedx.BOM)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, imageName, tag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSBOMDigestWithTag provides a mock function with given fields: ctx, repoName, tag
func (_m *GCRClient) GetSBOMDigestWithTag(ctx context.Context, repoName string, tag string) (string, error) {
	ret := _m.Called(ctx, repoName, tag)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, repoName, tag)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, repoName, tag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSBOMFromImage provides a mock function with given fields: sbom
func (_m *GCRClient) GetSBOMFromImage(sbom *v1.Image) (cyclonedx.BOM, error) {
	ret := _m.Called(sbom)

	var r0 cyclonedx.BOM
	if rf, ok := ret.Get(0).(func(*v1.Image) cyclonedx.BOM); ok {
		r0 = rf(sbom)
	} else {
		r0 = ret.Get(0).(cyclonedx.BOM)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1.Image) error); ok {
		r1 = rf(sbom)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVulnerabilities provides a mock function with given fields: ctx, ref
func (_m *GCRClient) GetVulnerabilities(ctx context.Context, ref reference.Canonical) ([]voucher.Vulnerability, error) {
	ret := _m.Called(ctx, ref)

	var r0 []voucher.Vulnerability
	if rf, ok := ret.Get(0).(func(context.Context, reference.Canonical) []voucher.Vulnerability); ok {
		r0 = rf(ctx, ref)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]voucher.Vulnerability)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, reference.Canonical) error); ok {
		r1 = rf(ctx, ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGCRClient creates a new instance of GCRClient. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewGCRClient(t testing.TB) *GCRClient {
	mock := &GCRClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
