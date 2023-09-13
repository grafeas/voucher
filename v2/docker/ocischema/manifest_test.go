package ocischema

import (
	"testing"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/ocischema"
	"github.com/docker/distribution/reference"
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

func TestToManifestList(t *testing.T) {
	newManifest := NewTestOCIManifestList()
	assert.Nil(t, newManifest)

	testRef, err := reference.ParseNamed("docker.io/library/docker@sha256:25a7feece7050334e8bd478dc9b6031c24db7fe81b2665abe690698ec52074f2")
	require.NoError(t, err)
	manifest, err := ToManifest(nil, testRef, newManifest)
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

func NewTestOCIManifestList() *manifestlist.DeserializedManifestList {
	// ManifestList based off of docker:latest image
	testManifestDescriptors := []manifestlist.ManifestDescriptor{
		// Matched Manifest
		{
			Platform: manifestlist.PlatformSpec{
				OS:           "linux",
				Architecture: "amd64",
			},
			Descriptor: distribution.Descriptor{
				Digest:    "sha256:bbc57559ea5f6d7359f53c92bdfd386df0b1b0384591a24b7a6cf40b77343a4a",
				Size:      3327,
				MediaType: "application/vnd.oci.image.manifest.v1+json",
			},
		},
		// Wrong arch
		{
			Platform: manifestlist.PlatformSpec{
				OS:           "linux",
				Architecture: "arm64",
				Variant:      "v8",
			},
			Descriptor: distribution.Descriptor{
				Digest:    "sha256:779ce156c5fb1a44c72252f2167ef492914727be3fc0abba6c4199414b383f10",
				Size:      3327,
				MediaType: "application/vnd.oci.image.manifest.v1+json",
			},
		},
		// Unknown
		{
			Platform: manifestlist.PlatformSpec{
				OS:           "unknown",
				Architecture: "unknown",
			},
			Descriptor: distribution.Descriptor{
				Digest:    "sha256:31b162dbcfef47a911a9c3d68e295c259b383bd6d552864240020c7f6b7ec847",
				Size:      840,
				MediaType: "application/vnd.oci.image.manifest.v1+json",
			},
		},
	}

	deserializedManifest, err := manifestlist.FromDescriptorsWithMediaType(
		testManifestDescriptors,
		manifestlist.OCISchemaVersion.MediaType,
	)
	if err != nil {
		panic("failed to generate new OCI manifest list")
	}

	return deserializedManifest
}
