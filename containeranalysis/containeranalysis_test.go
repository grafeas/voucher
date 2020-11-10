package containeranalysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Shopify/voucher"
)

const testImageName = "gcr.io/alpine/alpine@sha256:297524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c"

const testResourceAddress = "resourceUrl=\"https://" + testImageName + "\""

func TestGrafeasHelperFunctions(t *testing.T) {
	imageData, err := voucher.NewImageData(testImageName)
	require.NoError(t, err)
	assert.Equal(t, resourceURL(imageData), testResourceAddress)
}
