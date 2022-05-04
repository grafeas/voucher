package sbomgcr

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	gcr "github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/stretchr/testify/assert"
)

type mockGCRService struct {
	manifests map[string]gcr.ManifestInfo
	tag       name.Tag
}

func NewMockGCRService(sbomTag string) GCRService {
	manifests := map[string]gcr.ManifestInfo{
		"sha256:8d6f75268a5320cdd5473acb891ec60fc481cd84d0ee9b04be8569a974608d4c": {
			Tags:      []string{sbomTag},
			MediaType: "application/vnd.docker.distribution.manifest.v2+json",
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
	// image, err := tarball.ImageFromPath("fixtures/clouddo-ui.tar", &m.tag)
	// if err != nil {
	// 	return nil, err
	// }

	// // Hack: Change the layer manifestType to "application/vnd.dsse.envelope.v1+json"
	image, err := readOCIImage("fixtures/clouddo-sbom-oci")
	if err != nil {
		return nil, err
	}

	return image, nil
}

func readOCIImage(path string) (v1.Image, error) {
	// image read image from oci manifest
	imagePath, err := layout.FromPath("fixtures/clouddo-sbom-oci")
	if err != nil {
		return nil, err
	}

	imageIdex, err := imagePath.ImageIndex()
	if err != nil {
		return nil, err
	}

	indexManifest, err := imageIdex.IndexManifest()
	if err != nil {
		return nil, err
	}

	if len(indexManifest.Manifests) < 1 {
		return nil, fmt.Errorf("no manifests found in image manifest")
	}

	digest := indexManifest.Manifests[0].Digest

	image, err := imagePath.Image(digest)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func TestGetSBOM(t *testing.T) {
	// TODO: make this test more robust
	mockService := NewMockGCRService("sha256-551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f.att")
	client := NewClient(mockService)
	ctx := context.Background()

	boms, err := client.GetSBOM(ctx, "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui", "sha256-551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f.att")
	assert.NoError(t, err)
	isSBOM := strings.Contains(boms.Metadata.Component.Name, "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui")
	assert.True(t, isSBOM)
}
