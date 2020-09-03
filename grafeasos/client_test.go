package grafeasos

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/docker/distribution/reference"
	grafeaspb "github.com/grafeas/client-go/0.1.0"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Shopify/voucher"
)

var basePath string

func TestMain(m *testing.M) {
	flag.StringVar(&basePath, "grafeasos", "", "the base path to the grafeas instance to use for testing")
	flag.Parse()
	os.Exit(m.Run())
}

func getCanonicalRef(t *testing.T) reference.Canonical {
	named, err := reference.ParseNamed("us.gcr.io/grafeas/grafeas-server@sha256:c7303bdd6e36868d54b5b00dee125445a8d0f667c366420ccbe41dcf3b1c7733")
	require.NoError(t, err, "named")
	canonicalRef, err := reference.WithDigest(named, digest.FromString("sha256:c7303bdd6e36868d54b5b00dee125445a8d0f667c366420ccbe41dcf3b1c7733"))
	require.NoError(t, err, "canonicalRef")
	return canonicalRef
}

func newClientTest(t *testing.T, project string) voucher.MetadataClient {
	config := grafeaspb.NewConfiguration()
	if basePath != "" {
		config.BasePath = basePath
	} else {
		config.BasePath = "http://localhost:8080"
	}
	client, err := NewClient(context.Background(), project, project, nil, config)
	require.NoError(t, err)
	return client
}

func TestGrafeasMetadataClient(t *testing.T) {
	t.Run("Test CanAttest", func(t *testing.T) {
		client := newClientTest(t, "grafeasclienttest")

		canAttest := client.CanAttest()
		assert.False(t, canAttest)
	})

	t.Run("Test NewPayloadBody", func(t *testing.T) {
		client := newClientTest(t, "grafeasclienttest")
		ref := getCanonicalRef(t)

		_, err := client.NewPayloadBody(ref)
		require.NoError(t, err)
	})

	t.Run("Test AddAttestationToImage", func(t *testing.T) {
		ctx := context.Background()
		client := newClientTest(t, "grafeasclienttest")
		ref := getCanonicalRef(t)

		_, err := client.AddAttestationToImage(ctx, ref, voucher.Attestation{})
		require.Error(t, err)
	})

	if basePath != "" {
		PopulateGrafeasTestData(basePath)

		// t.Run("Test GetVulnerabilities with No Data", func(t *testing.T) {
		// 	ctx := context.Background()
		// 	client := newClientTest(t, "test")
		// 	ref := getCanonicalRef(t)
		// 	_, err := client.GetVulnerabilities(ctx, ref)
		// 	require.Error(t, err)
		// 	// require.Equal(t, err, errDiscoveriesUnfinished)
		// })

		t.Run("Test GetVulnerabilities with Data", func(t *testing.T) {
			// ctx := context.Background()
			// client := newClientTest(t, "grafeasclienttest")
			// ref := getCanonicalRef(t)

			// _, err := client.GetVulnerabilities(ctx, ref)
			// require.NoError(t, err)
			// require.Equal(t, 1, len(items))
		})

		// t.Run("Test GetBuildDetail with No Data", func(t *testing.T) {
		// 	ctx := context.Background()
		// 	client := newClientTest(t, "test")
		// 	ref := getCanonicalRef(t)

		// 	_, err := client.GetBuildDetail(ctx, ref)
		// 	require.Error(t, err)
		// 	_, ok := err.(*voucher.NoMetadataError)
		// 	require.True(t, ok)
		// })

		// 	t.Run("Test GetBuildDetail with Data", func(t *testing.T) {
		// 		ctx := context.Background()
		// 		client := newClientTest(t, "grafeasclienttest")
		// 		ref := getCanonicalRef(t)

		// 		_, err := client.GetBuildDetail(ctx, ref)
		// 		t.Error(t, "blal blee bloo")
		// 		require.NoError(t, err)
		// 	})
		// } else {
		// 	t.Skip("Test GetVulnerabilities with No Data")
		// 	t.Skip("Test GetVulnerabilities with Data")
		// 	t.Skip("Test GetBuildDetail with No Data")
		// 	t.Skip("Test GetBuildDetail with Data")
		// }
	}
}
