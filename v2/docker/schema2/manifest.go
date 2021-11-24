package schema2

import (
	"fmt"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	v2 "github.com/docker/distribution/manifest/schema2"
)

// IsManifest returns true if the passed manifest is a schema2 manifest.
func IsManifest(m distribution.Manifest) bool {
	switch m.(type) {
	case *v2.DeserializedManifest, *manifestlist.DeserializedManifestList:
		return true
	default:
		return false
	}
}

// ToManifest casts a distribution.Manifest to a schema2.Manifest. It panics
// if it passed anything other than a schema2.DeserialzedManifest.
func ToManifest(manifest distribution.Manifest) (v2.Manifest, error) {
	switch m := manifest.(type) {
	case *v2.DeserializedManifest:
		return m.Manifest, nil
	case *manifestlist.DeserializedManifestList:
		// TODO
		return v2.Manifest{}, fmt.Errorf("implement me")
	default:
		return v2.Manifest{}, fmt.Errorf("schema2.ToManifest was passed a %T", manifest)
	}
}
