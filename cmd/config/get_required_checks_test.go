package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequiredChecksFromConfig(t *testing.T) {
	FileName = "../../testdata/config.toml"
	InitConfig()

	expected := map[string][]string{
		"env1": {"diy"},
		"env2": {"diy", "nobody"},
	}
	got := GetRequiredChecksFromConfig()

	assert.Equal(t, len(expected), len(got))
	for groupName, expectedChecks := range expected {
		gotChecks, ok := got[groupName]
		assert.True(t, ok)
		assert.ElementsMatch(t, expectedChecks, gotChecks)
	}
}
