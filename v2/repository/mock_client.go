package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetCommit(ctx context.Context, details BuildDetail) (Commit, error) {
	args := m.Called(ctx, details)
	return args.Get(0).(Commit), args.Error(1)
}

func (m *MockClient) GetOrganization(ctx context.Context, details BuildDetail) (Organization, error) {
	args := m.Called(ctx, details)
	return args.Get(0).(Organization), args.Error(1)
}

func (m *MockClient) GetBranch(ctx context.Context, details BuildDetail, name string) (Branch, error) {
	args := m.Called(ctx, details)
	return args.Get(0).(Branch), args.Error(1)
}

func (m *MockClient) GetDefaultBranch(ctx context.Context, details BuildDetail) (Branch, error) {
	args := m.Called(ctx, details)
	return args.Get(0).(Branch), args.Error(1)
}
