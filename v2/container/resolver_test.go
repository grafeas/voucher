package container_test

import (
	"context"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/grafeas/voucher/v2/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolver_Resolve(t *testing.T) {
	t.Skip("depends on external tags which change over time - please update before running")

	googleAuth, err := google.NewEnvAuthenticator()
	require.NoError(t, err)
	cases := map[string]struct {
		ref      string
		expected string
		auth     authn.Authenticator
	}{
		"dockerhub multi-arch tag": {
			ref:      "debian:bullseye",
			expected: "index.docker.io/library/debian:bullseye@sha256:bfe6615d017d1eebe19f349669de58cda36c668ef916e618be78071513c690e5",
		},
		"dockerhub multi-arch tag+digest": {
			ref:      "debian:bullseye@sha256:bfe6615d017d1eebe19f349669de58cda36c668ef916e618be78071513c690e5",
			expected: "debian:bullseye@sha256:bfe6615d017d1eebe19f349669de58cda36c668ef916e618be78071513c690e5",
		},
		"dockerhub single-arch tag+digest": {
			ref:      "debian:bullseye@sha256:725ea075576e253aff7f205e8ff1f07d28ee28cb98089fb6d50eda013aeaeca5",
			expected: "debian:bullseye@sha256:725ea075576e253aff7f205e8ff1f07d28ee28cb98089fb6d50eda013aeaeca5",
		},
		"gcr latest": {
			ref:      "gcr.io/distroless/static-debian11",
			expected: "gcr.io/distroless/static-debian11:latest@sha256:5759d194607e472ff80fff5833442d3991dd89b219c96552837a2c8f74058617",
			auth:     googleAuth,
		},
		"gcr tag+digest": {
			ref:      "gcr.io/distroless/static-debian11:latest@sha256:5759d194607e472ff80fff5833442d3991dd89b219c96552837a2c8f74058617",
			expected: "gcr.io/distroless/static-debian11:latest@sha256:5759d194607e472ff80fff5833442d3991dd89b219c96552837a2c8f74058617",
			auth:     googleAuth,
		},
		"gcr tag": {
			ref:      "gcr.io/distroless/static-debian11:debug",
			expected: "gcr.io/distroless/static-debian11:debug@sha256:c66a6ecb5aa7704a68c89d3ead1398adc7f16e214dda5f5f8e5d44351bcbf67d",
			auth:     googleAuth,
		},
	}

	for label, tc := range cases {
		t.Run(label, func(t *testing.T) {
			res := container.NewResolver(tc.auth)
			resolved, err := res.ToDigest(context.Background(), tc.ref)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, resolved.String())
		})
	}
}
