package sbomgcr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCredentials(t *testing.T) {
	client, err := NewClient()
	man, _ := client.GetSBOM(context.Background(), "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui@sha256:551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f")
	fmt.Printf("%v\n", man)
	assert.NoError(t, err)
}
