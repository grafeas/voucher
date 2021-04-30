package schema2

import (
	"github.com/docker/distribution"
	v2 "github.com/docker/distribution/manifest/schema2"
)

// IsManifest returns true if the passed manifest is a schema2 manifest.
func IsManifest(m distribution.Manifest) bool {
	_, ok := m.(*v2.DeserializedManifest)
	return ok
}

// ToManifest casts a distribution.Manifest to a schema2.Manifest. It panics
// if it passed anything other than a schema2.DeserialzedManifest.
func ToManifest(manifest distribution.Manifest) v2.Manifest {
	schema2Manifest, ok := manifest.(*v2.DeserializedManifest)
	if !ok {
		panic("schema2.ToManifest was passed a non-schema2.DeserializedManifest")
	}

	return schema2Manifest.Manifest
}
