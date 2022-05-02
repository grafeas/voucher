package sbomgcr

import (
	"context"
	"fmt"
	"testing"

	"github.com/grafeas/voucher/v2/sbomgcr/mocks"
)

func TestGetSBOM(t *testing.T) {
	// TODO: make this test more robust
	service := mocks.NewGCRService(t)
	client := NewClient(service)
	ctx := context.Background()

	boms, _ := client.GetSBOM(ctx, "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui", "sha256-551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f.att")
	fmt.Println(boms)
}
