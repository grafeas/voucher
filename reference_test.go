package voucher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageReference(t *testing.T) {
	assert := assert.New(t)

	_, err := NewImageReference("!!abc")
	if assert.Error(err) {
		assert.Equal("can't use URL \"!!abc\" as image reference: invalid reference format", err.Error())
	}

	_, err = NewImageReference("gcr.io/path/to/image")
	if assert.Error(err) {
		assert.Equal("reference gcr.io/path/to/image has no digest", err.Error())
	}

	_, err = NewImageReference("gcr.io/path/to/image@sha256:97db2bc359ccc94d3b2d6f5daa4173e9e91c513b0dcd961408adbb95ec5e5ce5")
	assert.NoError(err)
}
