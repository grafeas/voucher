package sbomgcr

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSBOM(t *testing.T) {
	mockService := NewMockGCRService("sha256-551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f.att", "fixtures/clouddo-sbom-oci")
	client := NewClient(mockService)
	ctx := context.Background()

	boms, err := client.GetSBOM(ctx, "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui", "sha256-551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f.att")
	assert.NoError(t, err)
	isSBOM := strings.Contains(boms.Metadata.Component.Name, "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui")
	assert.True(t, isSBOM)
}
