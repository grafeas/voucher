package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	voucher "github.com/grafeas/voucher/v2"
)

func newScanner(secrets *Secrets, metadataClient voucher.MetadataClient, auth voucher.Auth) (scanner voucher.VulnerabilityScanner) {
	scannerName := viper.GetString("scanner")
	switch scannerName {
	case "gca", "g":
		log.Warningf("the %s option for `scanner` has been deprecated and will be removed in the future. Please use `metadata` instead.", scannerName)
		scanner = voucher.NewScanner(metadataClient)
	case "metadata":
		scanner = voucher.NewScanner(metadataClient)
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
