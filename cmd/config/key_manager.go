package config

import (
	"fmt"

	"github.com/Shopify/voucher"
	"github.com/spf13/viper"
)

func keyring() (*voucher.KeyRing, error) {
	if !viper.IsSet("ejson.dir") {
		return nil, fmt.Errorf("ESON dir not set in the config file")
	}

	dir := viper.GetString("ejson.dir")
	if !viper.IsSet("ejson.secrets") {
		return nil, fmt.Errorf("ESON secrets not set in the config file")
	}

	secrets := viper.GetString("ejson.secrets")

	return voucher.EjsonToKeyRing(dir, secrets)
}
