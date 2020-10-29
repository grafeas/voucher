package schema2

import (
	"testing"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vtesting "github.com/grafeas/voucher/testing"
)

func TestRequestConfig(t *testing.T) {
	ref := vtesting.NewTestReference(t)

	manifest := vtesting.NewTestManifest()

	client, server := vtesting.PrepareDockerTest(t, ref)
	defer server.Close()

	config, err := RequestConfig(client, ref, manifest)
	require.NoError(t, err, "failed to get config: %s", err)

	expectedConfig := &dockerTypes.ExecConfig{
		User: "nobody",
	}

	assert.Equal(t, expectedConfig, config)
}
