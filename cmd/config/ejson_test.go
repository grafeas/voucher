package config

import (
	"testing"

	"github.com/Shopify/voucher/clair"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNonExistantEjson(t *testing.T) {
	viper.Set("ejson.secrets", "../../tests/fixtures/bad.ejson")
	viper.Set("ejson.dir", "../../tests/fixtures/key")

	data := new(struct{})

	err := readEjson(data)
	if "stat ../../tests/fixtures/bad.ejson: no such file or directory" != err.Error() {
		t.Fatalf("did not fail appropriately, actual error is: %s", err)
	}
}

func TestGetClairConfig(t *testing.T) {
	viper.Set("ejson.secrets", "../../tests/fixtures/test.ejson")
	viper.Set("ejson.dir", "../../tests/fixtures/key")

	clairConfig, err := getClairConfig()
	assert.Nil(t, err)
	assert.NotNil(t, clairConfig)
	assert.Equal(t, clair.Config{
		Username: "testuser",
		Password: "testpassword",
	}, clairConfig)

}

func TestGetKeyRing(t *testing.T) {
	viper.Set("ejson.secrets", "../../tests/fixtures/test.ejson")
	viper.Set("ejson.dir", "../../tests/fixtures/key")

	keyRing, err := getKeyRing()
	assert.Nil(t, err)
	assert.NotNil(t, keyRing)

}
