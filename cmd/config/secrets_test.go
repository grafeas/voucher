package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafeas/voucher/clair"
	"github.com/grafeas/voucher/repository"
)

func TestNonExistantEjson(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/bad.ejson")
	viper.Set("ejson.dir", "../../testdata/key")

	_, err := ReadSecrets()
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

	data, err := ReadSecrets()
	require.NoError(t, err)
	assert.Equal(t, repository.KeyRing{
		"organization-name": repository.Auth{
			Token: "asdf1234",
		},
		"organization2-name": repository.Auth{
			Username: "testUser",
			Password: "testPassword",
		},
	}, data.RepositoryAuthentication)
}

func TestGetRepositoryKeyRingNoEjson(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test.ejson")
	viper.Set("ejson.dir", "../../testdata/nokey")

	data, err := ReadSecrets()
	require.Nil(t, data)
	assert.Error(t, err)
}

func TestGetClairConfig(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test.ejson")
	viper.Set("ejson.dir", "../../testdata/key")

	data, err := ReadSecrets()
	require.NoError(t, err)
	assert.Equal(t, clair.Config{
		Username: "testuser",
		Password: "testpassword",
	}, data.ClairConfig)
}

func TestGetPGPKeyRing(t *testing.T) {
	viper.Set("ejson.secrets", "../../testdata/test.ejson")
	viper.Set("ejson.dir", "../../testdata/key")

	data, err := ReadSecrets()
	require.NoError(t, err)
	keyRing, err := data.getPGPKeyRing()
	require.NoError(t, err)
	assert.NotNil(t, keyRing)
}
