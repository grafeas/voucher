package voucher

import (
	"context"
	"errors"
)

var errBrokenTest = errors.New("this test is broken")

// testCheck is a check that works and will pass depending on what the `shouldPass`
// variable is set to.
type testCheck struct {
	shouldPass bool
}

func (t *testCheck) Check(ctx context.Context, i ImageData) (bool, error) {
	return t.shouldPass, nil
}

// testBrokenCheck is a check that is completely broken and always returns an error.
type testBrokenCheck struct {
}

func (t *testBrokenCheck) Check(ctx context.Context, i ImageData) (bool, error) {
	return true, errBrokenTest
}

func newTestCheck(shouldPass bool) *testCheck {
	check := new(testCheck)
	check.shouldPass = shouldPass
	return check
}
