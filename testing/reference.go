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

// NewNobodyBadTestReference creates a new reference to be used for the nobody check.
// The returned reference is assumed to not, and does not have valid configuration
// or layers.
func NewNobodyBadTestReference(t *testing.T) reference.Canonical {
	t.Helper()

	return parseReference(t, "localhost/path/to/image@sha256:b248c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da")
}

// NewRateLimitedTestReference creates a new reference to be used to test
// the handling of Rate Limited docker calls.
// The returned response from calling this reference cannot be parsed by a
// JSON decoder and should result in an error.
func NewRateLimitedTestReference(t *testing.T) reference.Canonical {
	t.Helper()

	return parseReference(t, "localhost/path/to/ratelimited@sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da")
}

// NewTestSchema1Reference creates a new schema version 1 reference to be used
// throughout the docker tests. The returned reference is assumed to exist, and
// is assumed to have valid configuration and layers.
func NewTestSchema1Reference(t *testing.T) reference.Canonical {
	t.Helper()

	return parseReference(t, "localhost/schema1image@sha256:03f65aeeb2e8e8db022b297cae4cdce9248633f551452e63ba520d1f9ef2eca0")
}

// NewTestSchema1SignedReference creates a new schema version 1 reference to be
// used throughout the docker tests. The returned reference is assumed to
// exist, be signed, and to have valid configuration and layers.
func NewTestSchema1SignedReference(t *testing.T) reference.Canonical {
	t.Helper()

	return parseReference(t, "localhost/schema1imagesigned@sha256:18e6e7971438ab792d13563dcd8972acf4445bc0dcfdff84a6374d63a9c3ed62")
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
