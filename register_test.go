package voucher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCheckFactory(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	factories := make(CheckFactories)

	factories.Register("goodcheck", func() Check {
		return newTestCheck(true)
	})
	factories.Register("badcheck", func() Check {
		return newTestCheck(false)
	})

	checks, err := factories.GetNewChecks("goodcheck", "badcheck")
	assert.NoError(err)
	assert.Len(checks, 2)

	i := newTestImageData(t)

	if assert.NotNil(checks["goodcheck"]) {
		ok, checkErr := checks["goodcheck"].Check(ctx, i)
		assert.Nil(checkErr)
		assert.True(ok)
	}

	if assert.NotNil(checks["badcheck"]) {
		ok, checkErr := checks["badcheck"].Check(ctx, i)
		assert.Nil(checkErr)
		assert.False(ok)
	}
}

func TestEmptyCheckFactory(t *testing.T) {
	factories := make(CheckFactories)
	_, err := factories.GetNewChecks("nilcheck")
	assert.Contains(t, err.Error(), "requested check \"nilcheck\" does not exist")
}

func TestRegisterDefaultCheckFactories(t *testing.T) {
	assert := assert.New(t)

	// clear the existing CheckFactories, which should be empty regardless.
	DefaultCheckFactories = make(CheckFactories)

	RegisterCheckFactory("goodcheck", func() Check {
		return newTestCheck(true)
	})
	RegisterCheckFactory("badcheck", func() Check {
		return newTestCheck(false)
	})
	assert.Truef(IsCheckFactoryRegistered("goodcheck"), "goodcheck was registered but IsCheckRegistered is false")
	assert.Truef(IsCheckFactoryRegistered("badcheck"), "badcheck was registered but IsCheckRegistered is false")
	assert.False(IsCheckFactoryRegistered("nilcheck"), "nilcheck was not registered but IsCheckRegistered is true")

	checks, err := GetCheckFactories("nilcheck")
	assert.Error(err)
	assert.Equal(err.Error(), "requested check \"nilcheck\" does not exist")
	assert.Len(checks, 0)

	checks, err = GetCheckFactories("goodcheck", "badcheck")
	assert.NoError(err)
	assert.Len(checks, 2)

	// clear the existing CheckFactories, which should be empty regardless.
	DefaultCheckFactories = make(CheckFactories)
}
