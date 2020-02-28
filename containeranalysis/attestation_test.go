package containeranalysis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCheckNameFromNoteName(t *testing.T) {
	testValues := []struct {
		input    string
		expected string
	}{
		{
			input:    "projects/testproject/notes/diy",
			expected: "diy",
		},
		{
			input:    "projects/testproject/notes/",
			expected: "unknown",
		},
		{
			input:    "",
			expected: "unknown",
		},
	}

	for _, test := range testValues {
		output := getCheckNameFromNoteName("testproject", test.input)
		assert.Equal(t, test.expected, output)
	}
}
