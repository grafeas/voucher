package vtesting

import (
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/require"
)

// NewTestReference creates a new reference to be used throughout the docker tests.
// The returned reference is assumed to exist, and is assumed to have valid configuration
// and layers.
func NewTestReference(t *testing.T) reference.Canonical {
	t.Helper()

	return parseReference(t, "localhost/path/to/image@sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da")
}

// NewBadTestReference creates a new reference to be used throughout the docker tests.
// The returned reference is assumed to not, and does not have valid configuration
// or layers.
func NewBadTestReference(t *testing.T) reference.Canonical {
	t.Helper()

	return parseReference(t, "localhost/path/to/bad/image@sha256:bad8c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da")
}

// parseReference parses the passed reference and returns it (or fails)
func parseReference(t *testing.T, name string) reference.Canonical {
	t.Helper()

	ref, err := reference.Parse(name)
	require.NoErrorf(t, err, "could not make image reference (\"%s\"): %s", name, err)

	refCanonical, ok := ref.(reference.Canonical)
	require.True(t, ok, "could not convert reference to Canonical reference")

	return refCanonical

}
