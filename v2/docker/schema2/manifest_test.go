package schema2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vtesting "github.com/grafeas/voucher/v2/testing"
)

func TestToManifest(t *testing.T) {
	newManifest := vtesting.NewTestManifest()

	manifest, err := ToManifest(nil, nil, newManifest)
	require.NoError(t, err)
	assert.NotNil(t, manifest)
}
