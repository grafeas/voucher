package docker

import (
	"testing"

	"github.com/docker/distribution/reference"
)

const (
	testHostname    = "gcr.io"
	testProject     = "test/project"
	testDigest      = "sha256:cb749360c5198a55859a7f335de3cf4e2f64b60886a2098684a2f9c7ffca81f2"
	testBlobURL     = "https://" + testHostname + "/v2/" + testProject + "/blobs/" + testDigest
	testManifestURL = "https://" + testHostname + "/v2/" + testProject + "/manifests/" + testDigest
	testTokenURL    = "https://" + testHostname + "/v2/token?scope=repository%3Atest%2Fproject%3A%2A&service=gcr.io"
)

func compareStrings(t *testing.T, a, b string) {
	t.Helper()

	if a != b {
		t.Errorf("Passed Strings don't match:\n========\n%s\n========\n%s\n========\n", a, b)
	}
}

func TestGetBaseURI(t *testing.T) {
	named, err := reference.ParseNamed(testHostname + "/" + testProject + "@" + testDigest)
	if nil != err {
		t.Fatalf("failed to parse uri: %s", err)
	}

	compareStrings(t, testTokenURL, GetTokenURI(named))

	if canonicalRef, ok := named.(reference.Canonical); ok {
		compareStrings(t, string(canonicalRef.Digest()), testDigest)
		hostname, path := reference.SplitHostname(canonicalRef)
		compareStrings(t, hostname, "gcr.io")
		compareStrings(t, path, testProject)
		compareStrings(t, testBlobURL, GetBlobURI(canonicalRef, canonicalRef.Digest()))
		compareStrings(t, testManifestURL, GetManifestURI(canonicalRef))
	} else {
		t.Fatal("failed to get reference")
	}

}
