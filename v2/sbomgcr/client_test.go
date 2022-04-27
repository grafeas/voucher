package sbomgcr

// import (
// 	"context"

// 	"github.com/CycloneDX/cyclonedx-go"
// 	"github.com/docker/distribution/reference"
// )

// type MockClient struct{}

// func NewMockClient() *client {
// 	return &client{}
// }

// func (mc *MockClient) GetSBOM(ctx context.Context, ref reference.Canonical) (cyclonedx.BOM, error) {
// 	return cyclonedx.BOM{}, nil
// }

// func TestGetSBOM(t *testing.T) {
// 	// TODO:CS come back and actually make the mocks for these, don't merge this in! It's making network requests!
// 	client := NewClient()
// 	img := "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui@sha256:551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f", "sha256:551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f"
// 	ref := getCanonicalRef(t, img, digest)
// 	man, _ := client.GetSBOM(context.Background(), img)
// 	fmt.Printf("%v\n", man)
// }

// func TestGetSBOMDigestWithTag(t *testing.T) {
// 	client := NewClient()
// 	img, digest := "gcr.io/shopify-codelab-and-demos/sbom-lab/apps/production/clouddo-ui@sha256:551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f", "sha256:551182244aa6ab6997900bc04dd4e170ef13455c068360e93fc7b149eb2bc45f"
// 	ref := getCanonicalRef(t, img, digest)
// 	sbomTag := GetSBOMTagFromImage(ref)

// 	digest, err := client.GetSBOMDigestWithTag(context.Background(), ref.Name(), sbomTag)
// 	require.NoError(t, err, "digest")
// }

// func getCanonicalRef(t *testing.T, img string, digestStr string) reference.Canonical {
// 	named, err := reference.ParseNamed(img)
// 	require.NoError(t, err, "named")
// 	canonicalRef, err := reference.WithDigest(named, digest.Digest(digestStr))
// 	require.NoError(t, err, "canonicalRef")
// 	return canonicalRef
// }
