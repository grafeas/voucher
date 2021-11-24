package schema2

import (
	"fmt"
	"net/http"
	"os"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	v2 "github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
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
func ToManifest(client *http.Client, ref reference.Named, manifest distribution.Manifest) (v2.Manifest, error) {
	switch m := manifest.(type) {
	case *v2.DeserializedManifest:
		return m.Manifest, nil
	case *manifestlist.DeserializedManifestList:
		return resolveManifestFromList(client, ref, m)
	default:
		return v2.Manifest{}, fmt.Errorf("schema2.ToManifest was passed a %T", manifest)
	}
}

// Ugly method to override the target os/arch without wiring the voucher config to this context
var targetOS, targetArch string

func init() {
	targetOS = os.Getenv("VOUCHER_TARGET_OS")
	if targetOS == "" {
		targetOS = "linux"
	}
	targetArch = os.Getenv("VOUCHER_TARGET_ARCH")
	if targetArch == "" {
		targetArch = "amd64"
	}
}

func resolveManifestFromList(client *http.Client, ref reference.Named, mfs *manifestlist.DeserializedManifestList) (v2.Manifest, error) {
	for _, mf := range mfs.Manifests {
		if mf.Platform.Architecture != targetArch || mf.Platform.OS != targetOS {
			continue
		}

		// TODO: use the client and the fancy new uri.Manifest helper to grab this platform's manifest, return that
		// that's doable - this commit is focussed on wiring dependencies to this function
	}
	return v2.Manifest{}, fmt.Errorf("no manifest matching %s/%s found", targetOS, targetArch)
}
