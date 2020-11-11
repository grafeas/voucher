package schema1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vtesting "github.com/grafeas/voucher/testing"
)

func TestConfigFromManifest(t *testing.T) {
	pk := vtesting.NewPrivateKey()
	newManifest := vtesting.NewTestSchema1SignedManifest(pk)

	// we can pass nil as the http.Client because schema1's config is stored in
	// the history fields. It's super weird.
	config, err := RequestConfig(nil, nil, newManifest)
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "nobody", config.User)
}
