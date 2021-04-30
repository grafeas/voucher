package uri

import (
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testImageReference = "gcr.io/alpine/alpine@sha256:297524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c"

func TestReferenceToProjectName(t *testing.T) {
	ref, err := reference.Parse(testImageReference)
	require.NoError(t, err)

	project, err := ReferenceToProjectName(ref)
	assert.NoError(t, err)
	assert.Equal(t, "alpine", project)
}

func TestReferenceToProjectNameWithOtherReference(t *testing.T) {
	ref, err := reference.Parse("alpine/alpine")

	require.NoError(t, err)

	project, err := ReferenceToProjectName(ref)
	assert.Error(t, err)
	assert.Equal(t, "", project)
}
