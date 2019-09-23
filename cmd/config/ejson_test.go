package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Shopify/voucher/clair"
	"github.com/Shopify/voucher/repository"
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

func TestGetRepositoryKeyRing(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test.ejson")
	viper.Set("ejson.dir", "../../testdata/key")

	repoKeyRing, err := getRepositoryKeyRing()
	require.NoError(t, err)
	assert.NotNil(t, repoKeyRing)
	assert.Equal(t, repository.KeyRing{
		"organization-name": repository.Auth{
			Token: "asdf1234",
		},
		"organization2-name": repository.Auth{
			Username: "testUser",
			Password: "testPassword",
		},
	}, repoKeyRing)
}

func TestGetRepositoryKeyRingNoEjson(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test.ejson")
	viper.Set("ejson.dir", "../../testdata/nokey")

	repoKeyRing, err := getRepositoryKeyRing()
	require.Nil(t, repoKeyRing)
	assert.Error(t, err)
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
