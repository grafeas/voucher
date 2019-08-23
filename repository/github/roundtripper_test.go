package github

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoundTripperWrapper(t *testing.T) {
	client := &http.Client{}
	rtw := newRoundTripperWrapper(client.Transport)

	assert.Implements(t, (*http.RoundTripper)(nil), rtw)
}

func TestAddPreviewSchemaHeaders(t *testing.T) {
	headerTests := []struct {
		testName string
		headers  []string
		expected http.Header
	}{
		{
			testName: "Headers exist",
			headers:  []string{"header1", "header2"},
			expected: map[string][]string{
				"Accept": []string{"header1", "header2"},
			},
		},
		{
			testName: "Headers do not exist",
			headers:  []string{},
			expected: map[string][]string{},
		},
		{
			testName: "Parameters do not exist",
			headers:  nil,
			expected: map[string][]string{},
		},
	}

	for _, test := range headerTests {
		t.Run(test.testName, func(t *testing.T) {
			req, err := http.NewRequest("GET", "https://github.com/", nil)

			assert.NoError(t, err)

			newReq := addPreviewSchemaHeaders(req, test.headers)
			assert.Exactly(t, newReq.Header, test.expected)
		})
	}
}

func TestRoundTrip(t *testing.T) {
	client := &http.Client{
		Transport: new(http.Transport),
	}
	rtw := newRoundTripperWrapper(client.Transport)
	req, err := http.NewRequest("GET", "https://github.com/", nil)

	assert.NoError(t, err)

	res, err := rtw.RoundTrip(req)

	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
}
