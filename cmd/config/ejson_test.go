package config

import (
	"testing"

	"github.com/Shopify/voucher/clair"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNonExistantEjson(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/bad.ejson")
	viper.Set("ejson.dir", "../../testdata/key")

	data := new(struct{})

	err := readEjson(data)
	require.Equal(
		t,
		err.Error(),
		"stat ../../testdata/bad.ejson: no such file or directory",
		"did not fail appropriately, actual error is:",
		err,
	)
}

func TestGetClairConfig(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test.ejson")
	viper.Set("ejson.dir", "../../testdata/key")

	clairConfig, err := getClairConfig()
	require.NoError(t, err)
	require.NotNil(t, clairConfig)
	assert.Equal(t, clair.Config{
		Username: "testuser",
		Password: "testpassword",
	}, clairConfig)
}

func TestGetKeyRing(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test.ejson")
	viper.Set("ejson.dir", "../../testdata/key")

	keyRing, err := getKeyRing()
	require.NoError(t, err)
	assert.NotNil(t, keyRing)
}
