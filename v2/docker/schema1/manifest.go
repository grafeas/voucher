package schema1

import (
	"github.com/docker/distribution"
	v1 "github.com/docker/distribution/manifest/schema1"
)

// IsManifest returns true if the passed manifest is a schema1 manifest.
func IsManifest(m distribution.Manifest) bool {
	_, ok := m.(*v1.SignedManifest)
	return ok
}

// ToManifest casts a distribution.Manifest to a schema1.Manifest. It panics
// if it passed anything other than a schema1.SignedManifest.
func ToManifest(manifest distribution.Manifest) *v1.SignedManifest {
	signedManifest, ok := manifest.(*v1.SignedManifest)
	if !ok {
		panic("schema1.ToManifest was passed a non-schema1.SignedManifest")
	}

	return signedManifest
}
