package ocischema

import (
	"testing"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/ocischema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// MediaTypeOCILayer is the media type for OCI layers.
	OCILayer = "application/vnd.oci.image.layer.v1.tar+gzip"
)

func TestToManifest(t *testing.T) {
	newManifest := NewTestOCIManifest()
	manifest, err := ToManifest(nil, nil, newManifest)
	require.NoError(t, err)
	assert.NotNil(t, manifest)
}

func NewTestOCIManifest() *ocischema.DeserializedManifest {
	manifest := ocischema.Manifest{
		Config: distribution.Descriptor{},
		Layers: []distribution.Descriptor{
			{
				MediaType: OCILayer,
				Digest:    "sha256:31e352740f534f9ad170f75378a84fe453d6156e40700b882d737a8f4a6988a3",
				Size:      3397879,
			},
			{
				MediaType: OCILayer,
				Digest:    "sha256:a909a6dccb2a6d84f4f5aefc708bf3ffcc5d43ef0281708c595f1d3f126d395d",
				Size:      2014277,
			},
			{
				MediaType: OCILayer,
				Digest:    "sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1",
				Size:      32,
			},
		},
	}

	manifest.SchemaVersion = 2
	manifest.MediaType = ocischema.SchemaVersion.MediaType
	manifest.Config.Digest = "sha256:25a7feece7050334e8bd478dc9b6031c24db7fe81b2665abe690698ec52074f2"

	newManifest, err := ocischema.FromStruct(manifest)
	if err != nil {
		panic("failed to generate new OCI manifest")
	}
	return newManifest
}
