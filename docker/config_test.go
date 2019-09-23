package docker

import (
	"testing"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vtesting "github.com/Shopify/voucher/testing"
)

func TestRequestConfig(t *testing.T) {
	ref := vtesting.NewTestReference(t)

	client, server := PrepareDockerTest(t, ref)
	defer server.Close()

	config, err := RequestImageConfig(client, ref)
	require.NoError(t, err, "failed to get config: %s", err)

	expectedConfig := ImageConfig{
		ContainerConfig: dockerTypes.ExecConfig{
			User: "root",
		},
	}

	assert.Equal(t, expectedConfig, config)

	assert.True(t, config.RunsAsRoot())
}
