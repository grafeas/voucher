package voucher

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCheck struct {
	mock.Mock
}

func (m *MockCheck) Check(ctx context.Context, i ImageData) (bool, error) {
	args := m.Called(ctx, i)
	return args.Bool(0), args.Error(1)
}
