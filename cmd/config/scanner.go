package config

import (
	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/clair"
	"github.com/Shopify/voucher/grafeas"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func newScanner(metadataClient voucher.MetadataClient, auth voucher.Auth) (scanner voucher.VulnerabilityScanner) {
	scannerName := viper.GetString("scanner")
	switch scannerName {
	case "clair", "c":
		config, err := getClairConfig()
		if nil != err {
			log.Fatalf("could not load clair configuration: %s", err)
		}
		if "" == config.Hostname {
			config.Hostname = viper.GetString("clair.address")
		}
		scanner = clair.NewScanner(config, auth)
	case "gca", "g":
		scanner = grafeas.NewScanner(metadataClient)
	default:
		scanner = nil
	}

	if nil == scanner {
		log.Fatalf("not a valid scanner: %s", scannerName)
	}

	severity, err := voucher.StringToSeverity(viper.GetString("failon"))
	if nil != err {
		log.Fatal(err)
	}

	scanner.FailOn(severity)

	return
}
