package ocischema

import (
	"testing"

	vtesting "github.com/grafeas/voucher/v2/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToManifest(t *testing.T) {
	newManifest := vtesting.NewTestOCIManifest()
	manifest, err := ToManifest(nil, nil, newManifest)
	require.NoError(t, err)
	assert.NotNil(t, manifest)
}
