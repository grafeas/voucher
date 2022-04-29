package sbomgcr

import (
	"context"
	"fmt"
	"testing"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/grafeas/voucher/v2/sbomgcr/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSBOM(t *testing.T) {
	mockClient := &mocks.GCRClient{}
	mockClient.On("GetSBOM", context.Background(), mock.Anything, mock.Anything).Return(cyclonedx.BOM{}, nil)
	tests := map[string]struct {
		imageName   string
		tag         string
		expectedErr error
	}{
		"success": {
			imageName:   "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui",
			tag:         "sha256-inhere.att",
			expectedErr: nil,
		},
		"error on getting sbom digest with tag": {
			imageName:   "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui",
			tag:         "sha256-nothere.att",
			expectedErr: fmt.Errorf("error getting digest with tag no digest found in Client.GetSBOMDigestWithTag"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := mockClient.GetSBOM(context.Background(), test.imageName, test.tag)
			if err != nil {
				assert.Equal(t, test.expectedErr, err)
			}
			assert.NoError(t, err)
		})
	}
}
