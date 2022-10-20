package config

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/containeranalysis"
	"github.com/grafeas/voucher/v2/grafeas"
	"github.com/grafeas/voucher/v2/signer"
)

// NewMetadataClient creates a new MetadataClient.
func NewMetadataClient(ctx context.Context, secrets *Secrets) (voucher.MetadataClient, error) {
	keyring := NewAttestationSigner(secrets)

	if viper.GetString("image_project") != "" {
		log.Warning("`image_project` is deprecated. Please rely on the `valid_repos` configuration option to limit where images come from.")
	}

	metadataClient := viper.GetString("metadata_client")
	switch metadataClient {
	case "containeranalysis":
		return containeranalysis.NewClient(
			ctx,
			viper.GetString("binauth_project"),
			viper.GetString("containeranalysis.build_detail_fallback_project"),
			keyring,
		)
	case "grafeasos":
		return grafeas.NewClient(
			ctx,
			viper.GetString("binauth_project"),
			viper.GetString("grafeasos.vul_project"),
			keyring,
			grafeas.NewAPIService(viper.GetString("grafeasos.hostname"), viper.GetString("grafeasos.version")),
		)
	default:
		log.Warning("`metadata_client` option is not set, defaulting to \"containeranalysis\"")
		return containeranalysis.NewClient(
			ctx,
			viper.GetString("binauth_project"),
			viper.GetString("containeranalysis.build_detail_fallback_project"),
			keyring,
		)
	}
}

// NewAttestationSigner creates a new attestation signer
func NewAttestationSigner(secrets *Secrets) signer.AttestationSigner {
	signerName := viper.GetString("signer")
	if signerName == "pgp" || signerName == "" {
		if secrets == nil {
			log.Println("could not load PGP keyring from ejson - no secrets configured")
			return nil
		}
		keyring, err := secrets.getPGPKeyRing()
		if nil != err {
			log.Println("could not load PGP keyring from ejson, continuing without attestation support: ", err)
			return nil
		}
		return keyring
	} else if signerName == "kms" {
		keyring, err := getKMSKeyRing()
		if nil != err {
			log.Println("could not load KMS keyring from config, continuing without attestation support: ", err)
			return nil
		}
		return keyring
	}
	log.Printf("signer %q is unknown, supported values are 'kms' or 'pgp'\n", signerName)
	return nil
}
