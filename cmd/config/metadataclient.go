package config

import (
	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/grafeas"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewMetadataClient creates a new MetadataClient.
func NewMetadataClient() voucher.MetadataClient {
	keyring, err := keyring()
	if nil != err {
		log.Println("could not load keyring from ejson, continuing without attestation support: ", err)
		keyring = nil
	}

	return grafeas.NewClient(
		viper.GetString("image_project"),
		viper.GetString("binauth_project"),
		keyring,
	)
}
