package sbomgcr

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	gcr "github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

type mockGCRService struct {
	manifests map[string]gcr.ManifestInfo
	tag       name.Tag
}

func NewMockGCRService(sbomTag string) GCRService {
	manifests := map[string]gcr.ManifestInfo{
		"sha256:8d6f75268a5320cdd5473acb891ec60fc481cd84d0ee9b04be8569a974608d4c": {
			Tags: []string{sbomTag},
		},
	}
	tag, _ := name.NewTag("gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui:i-was-a-digest")
	return &mockGCRService{manifests: manifests, tag: tag}
}

func (m mockGCRService) ListTags(ctx context.Context, repo name.Repository) (*gcr.Tags, error) {
	return &gcr.Tags{
		Manifests: m.manifests,
	}, nil
}

func (m mockGCRService) PullImage(src string) (v1.Image, error) {
	img, err := tarball.ImageFromPath("fixtures/clouddo-ui.tar", &m.tag)
	fmt.Println(img, err)
	return img, err
}

func TestGetSBOM(t *testing.T) {
	// mockService := NewMockGCRService("sha256-551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f.att")
	// client := NewClient(mockService)
	// ctx := context.Background()

	// boms, err := client.GetSBOM(ctx, "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui", "sha256-551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f.att")
	// assert.NoError(t, err)
	fmt.Println("hello")
}
