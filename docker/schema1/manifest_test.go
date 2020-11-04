package schema1

import (
	"testing"

	"github.com/stretchr/testify/assert"

	vtesting "github.com/grafeas/voucher/testing"
)

func TestToManifest(t *testing.T) {
	pk := vtesting.NewPrivateKey()
	newManifest := vtesting.NewTestSchema1SignedManifest(pk)

	manifest := ToManifest(newManifest)
	assert.NotNil(t, manifest)
}
