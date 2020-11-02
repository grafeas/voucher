package config

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/grafeas/containeranalysis"
	"github.com/Shopify/voucher/grafeas/rest"
	"github.com/Shopify/voucher/signer"
)

// NewMetadataClient creates a new MetadataClient.
func NewMetadataClient(ctx context.Context, secrets *Secrets) (voucher.MetadataClient, error) {
	var keyring signer.AttestationSigner
	var err error

	signerName := viper.GetString("signer")
	if signerName == "pgp" || signerName == "" {
		if secrets == nil {
			log.Println("could not load PGP keyring from ejson - no secrets configured")
			keyring = nil
		} else {
			keyring, err = secrets.getPGPKeyRing()
			if nil != err {
				log.Println("could not load PGP keyring from ejson, continuing without attestation support: ", err)
				keyring = nil
			}
		}
	} else if signerName == "kms" {
		keyring, err = getKMSKeyRing()
		if nil != err {
			log.Println("could not load KMS keyring from config, continuing without attestation support: ", err)
			keyring = nil
		}
	} else {
		log.Printf("signer %q is unknown, supported values are 'kms' or 'pgp'\n", signerName)
	}

	if viper.GetString("image_project") != "" {
		log.Warning("`image_project` is deprecated. Please rely on the `valid_repos` configuration option to limit where images come from.")
	}

	metadataClient := viper.GetString("metadata_client")
	switch metadataClient {
	case "containeranalysis":
		return containeranalysis.NewClient(
			ctx,
			viper.GetString("binauth_project"),
			keyring,
		)
	case "grafeasos":
		return rest.NewClient(
			ctx,
			viper.GetString("image_project"),
			viper.GetString("binauth_project"),
			viper.GetString("grafeasos_vul_project"),
			keyring,
			rest.NewGrafeasAPIService(viper.GetString("grafeasos_base_path"), viper.GetString("grafeasos_version")),
		)
	default:
		log.Warning("`metadata_client` option is not set, defaulting to \"containeranalysis\"")
		return containeranalysis.NewClient(
			ctx,
			viper.GetString("binauth_project"),
			keyring,
		)
	}
}
