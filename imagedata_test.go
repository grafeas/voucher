package voucher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageData(t *testing.T) {
	assert := assert.New(t)

	_, err := NewImageData("!!abc")
	assert.NotNil(err)
	assert.Equal("can't use URL in ImageData: invalid reference format", err.Error())

	_, err = NewImageData("gcr.io/path/to/image")
	assert.NotNil(err)
	assert.Equal("reference gcr.io/path/to/image has no digest", err.Error())

	_, err = NewImageData("gcr.io/path/to/image@sha256:97db2bc359ccc94d3b2d6f5daa4173e9e91c513b0dcd961408adbb95ec5e5ce5")
	assert.Nil(err)
}
