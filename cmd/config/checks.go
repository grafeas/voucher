package config

import (
	// Register the DIY check
	_ "github.com/Shopify/voucher/checks/diy"

	// Register the Nobody check
	_ "github.com/Shopify/voucher/checks/nobody"

	// Register the Provenance check
	_ "github.com/Shopify/voucher/checks/provenance"

	// Register the Snakeoil check
	_ "github.com/Shopify/voucher/checks/snakeoil"
)

import (
	"fmt"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/clair"
	"github.com/Shopify/voucher/grafeas"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func newScanner(metadataClient voucher.MetadataClient) (scanner voucher.VulnerabilityScanner) {
	scannerName := viper.GetString("scanner")
	switch scannerName {
	case "clair", "c":
		scanner = clair.NewScanner(viper.GetString("clair.address"))
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

// EnabledChecks returns a slice of strings with the check names, based on a
// map[string]bool (with a check name in the key, and the value storing whether
// or not to run the check). The returned map contains enabled checks.
func EnabledChecks(checks map[string]bool) (enabledChecks []string) {
	enabledChecks = make([]string, 0, len(checks))
	for name, check := range checks {
		if check {
			enabledChecks = append(enabledChecks, name)
		}
	}
	return
}

// NewCheckSuite creates a new checks.Suite with the requested
// Checks, passing any necessary configuration details to the
// checks.
func NewCheckSuite(metadataClient voucher.MetadataClient, names ...string) (*voucher.Suite, error) {
	scanner := newScanner(metadataClient)
	checksuite := voucher.NewSuite()

	checks, err := voucher.GetCheckFactories(names...)
	if nil != err {
		return checksuite, fmt.Errorf("can't create check suite: %s", err)
	}

	for name, check := range checks {
		if vulCheck, ok := check.(voucher.VulnerabilityCheck); ok {
			vulCheck.SetScanner(scanner)
			checksuite.Add(name, vulCheck)
			continue
		}

		if metadataCheck, ok := check.(voucher.MetadataCheck); ok {
			metadataCheck.SetMetadataClient(metadataClient)
			checksuite.Add(name, metadataCheck)
			continue
		}

		checksuite.Add(name, check)
	}

	return checksuite, nil
}
