package gcr_test

import (
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/grafeas/voucher/v2/container/gcr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReferenceToProjectName(t *testing.T) {
	// map of reference to project
	cases := map[string]string{
		"gcr.io/alpine/alpine@sha256:297524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c": "alpine",
		"alpine/alpine": "",
		"southamerica-east1-docker.pkg.dev/my-project/team1/webapp":   "my-project",
		"australia-southeast1-docker.pkg.dev/my-project/team2/webapp": "my-project",
	}
	for img, expectedProject := range cases {
		t.Run(img, func(t *testing.T) {
			ref, err := reference.Parse(img)
			require.NoError(t, err)

			project, err := gcr.ReferenceToProjectName(ref)
			if expectedProject != "" {
				assert.NoError(t, err)
				assert.Equal(t, expectedProject, project)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
