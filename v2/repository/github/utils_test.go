package github

import (
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewGitHubV4URI(t *testing.T) {
	testURLs := []struct {
		testName string
		url      string
		expected string
	}{
		{
			testName: "Valid URL",
			url:      "https://www.github.com/commit/2389429034",
			expected: "https://www.github.com/commit/2389429034",
		},
		{
			testName: "Invalid URL",
			url:      "oof%^&*%*(",
			expected: "",
		},
	}

	for _, test := range testURLs {
		t.Run(test.testName, func(t *testing.T) {
			actualURL, err := createNewGitHubV4URI(test.url)

			if test.expected == "" {
				assert.Error(t, err, "Test did not throw an error when it should have")
				assert.Nil(t, actualURL, "Test returned a URL when it shouldn't have")
			} else {
				assert.NoError(t, err, "Test threw an error when it shouldn't have")
				assert.Equal(t, test.expected, actualURL.String(), "Test did not return a URL when it should have")
				assert.IsType(t, &githubv4.URI{}, actualURL, "Output should be of type githubv4.URI")
			}
		})
	}
}

func TestIsValidCheckConclusionState(t *testing.T) {
	stateTests := []struct {
		testName string
		state    checkConclusionState
		expected bool
	}{
		{
			testName: "Valid non-empty state",
			state:    "ACTION_REQUIRED",
			expected: true,
		},
		{
			testName: "Valid empty state",
			state:    "",
			expected: true,
		},
		{
			testName: "Invalid state",
			state:    "invalid",
			expected: false,
		},
	}

	for _, test := range stateTests {
		t.Run(test.testName, func(t *testing.T) {
			res := test.state.isValidCheckConclusionState()
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestIsValidCheckStatusState(t *testing.T) {
	stateTests := []struct {
		testName string
		state    checkStatusState
		expected bool
	}{
		{
			testName: "Valid non-empty state",
			state:    "QUEUED",
			expected: true,
		},
		{
			testName: "Valid empty state",
			state:    "",
			expected: true,
		},
		{
			testName: "Invalid state",
			state:    "invalid",
			expected: false,
		},
	}

	for _, test := range stateTests {
		t.Run(test.testName, func(t *testing.T) {
			res := test.state.isValidCheckStatusState()
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestIsValidStatusState(t *testing.T) {
	stateTests := []struct {
		testName string
		state    statusState
		expected bool
	}{
		{
			testName: "Valid non-empty state",
			state:    "PENDING",
			expected: true,
		},
		{
			testName: "Valid empty state",
			state:    "",
			expected: true,
		},
		{
			testName: "Invalid state",
			state:    "invalid",
			expected: false,
		},
	}

	for _, test := range stateTests {
		t.Run(test.testName, func(t *testing.T) {
			res := test.state.isValidStatusState()
			assert.Equal(t, test.expected, res)
		})
	}
}
