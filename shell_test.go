package voucher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunShellCommand(t *testing.T) {
	assert := assert.New(t)
	result, err := RunShellCommand("echo", "Hello World")
	assert.Nil(err)
	assert.Equal("Hello World\n", result)
}
