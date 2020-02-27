package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testExpectedSlice = []string{
	"a",
	"f",
}

var testGoodMap = map[string]interface{}{
	"a": true,
	"b": false,
	"c": false,
	"e": 55,
	"f": true,
}

func TestGetRequiredChecksFromConfig(t *testing.T) {
	FileName = "../../testdata/config.toml"
	InitConfig()

	expected := map[string][]string{
		"all": {
			"diy",
			"nobody",
			"provenance",
			"snakeoil",
		},
		"env1": {
			"diy",
		},
		"env2": {
			"diy",
			"nobody",
		},
	}
	got := GetRequiredChecksFromConfig()
	for groupName, expectedChecks := range expected {
		t.Logf("%s: %s", groupName, strings.Join(expectedChecks, ", "))
	}
	for groupName, expectedChecks := range got {
		t.Logf("%s: %s", groupName, strings.Join(expectedChecks, ", "))
	}

	assert.Equal(t, len(expected), len(got))
	for groupName, expectedChecks := range expected {
		gotChecks, ok := got[groupName]
		assert.True(t, ok)
		assert.ElementsMatch(t, expectedChecks, gotChecks)
	}
}

func TestToStringSlice(t *testing.T) {
	convert := toStringSlice(testGoodMap)
	assert.ElementsMatch(t, testExpectedSlice, convert)
}
