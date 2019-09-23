package config

import (
	"fmt"

	"github.com/Shopify/voucher"
	// Register the DIY check
	_ "github.com/Shopify/voucher/checks/diy"
	// Register the Nobody check
	_ "github.com/Shopify/voucher/checks/nobody"
	// Register the Provenance check
	_ "github.com/Shopify/voucher/checks/provenance"
	// Register the Snakeoil check
	_ "github.com/Shopify/voucher/checks/snakeoil"
)

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

// setAuth sets the Auth for the passed Check, if that Check implements
// AuthorizedCheck.
func setCheckAuth(check voucher.Check, auth voucher.Auth) {
	if authCheck, ok := check.(voucher.AuthorizedCheck); ok {
		authCheck.SetAuth(auth)
	}
}

// setCheckScanner sets the scanner on the passed Check, if that Check implements
// VulnerabilityCheck.
func setCheckScanner(check voucher.Check, scanner voucher.VulnerabilityScanner) {
	if vulCheck, ok := check.(voucher.VulnerabilityCheck); ok {
		vulCheck.SetScanner(scanner)
	}
}

// setCheckMetadataClient sets the MetadataClient for the passed Check, if that Check implements
// MetadataCheck.
func setCheckMetadataClient(check voucher.Check, metadataClient voucher.MetadataClient) {
	if metadataCheck, ok := check.(voucher.MetadataCheck); ok {
		metadataCheck.SetMetadataClient(metadataClient)
	}
}

// setCheckValidRepos sets the valid repos list for the passed Check, if
// that Check is a RepoValidatorCheck.
func setCheckValidRepos(check voucher.Check, validRepos []string) {
	if validRepoCheck, ok := check.(voucher.RepoValidatorCheck); ok {
		validRepoCheck.SetValidRepos(validRepos)
	}
}

// NewCheckSuite creates a new checks.Suite with the requested
// Checks, passing any necessary configuration details to the
// checks.
func NewCheckSuite(metadataClient voucher.MetadataClient, names ...string) (*voucher.Suite, error) {
	auth := newAuth()
	repos := validRepos()
	scanner := newScanner(metadataClient, auth)
	checksuite := voucher.NewSuite()

	checks, err := voucher.GetCheckFactories(names...)
	if nil != err {
		return checksuite, fmt.Errorf("can't create check suite: %s", err)
	}

	for name, check := range checks {
		setCheckAuth(check, auth)
		setCheckScanner(check, scanner)
		setCheckMetadataClient(check, metadataClient)
		setCheckValidRepos(check, repos)

		checksuite.Add(name, check)
	}

	return checksuite, nil
}
