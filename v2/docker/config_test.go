package docker

import (
	"testing"

	"github.com/stretchr/testify/require"

	vtesting "github.com/grafeas/voucher/v2/testing"
)

func TestRequestConfig(t *testing.T) {
	ref := vtesting.NewTestReference(t)

	client, server := vtesting.PrepareDockerTest(t, ref)
	defer server.Close()

	config, err := RequestImageConfig(client, ref)
	require.NoError(t, err)
	require.False(t, config.RunsAsRoot())
}
