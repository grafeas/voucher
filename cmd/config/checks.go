package config

import (
	"fmt"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/repository"
	"github.com/spf13/viper"

	// Register the DIY check
	_ "github.com/grafeas/voucher/checks/diy"
	// Register the Nobody check
	_ "github.com/grafeas/voucher/checks/nobody"
	// Register the Provenance check
	_ "github.com/grafeas/voucher/checks/provenance"
	// Register the Snakeoil check
	_ "github.com/grafeas/voucher/checks/snakeoil"
	// Register the Repo check
	_ "github.com/grafeas/voucher/checks/approved"
)

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

// setCheckTrustedIdentitiesAndProjects sets trusted identities and projects used in ProvenanceCheck
// trustedBuildCreators is a list of trusted build creators
// trustedProjects is a list of trusted projects
func setCheckTrustedIdentitiesAndProjects(check voucher.Check, trustedBuildCreators []string, trustedProjects []string) {
	if provenanceCheck, ok := check.(voucher.ProvenanceCheck); ok {
		provenanceCheck.SetTrustedBuildCreators(trustedBuildCreators)
		provenanceCheck.SetTrustedProjects(trustedProjects)
	}
}

// setCheckRepositoryClient sets the repository client for the passed Check, if that Check implements
// RepositoryCheck.
func setCheckRepositoryClient(check voucher.Check, repositoryClient repository.Client) {
	if repositoryCheck, ok := check.(voucher.RepositoryCheck); ok {
		repositoryCheck.SetRepositoryClient(repositoryClient)
	}
}

// NewCheckSuite creates a new checks.Suite with the requested
// Checks, passing any necessary configuration details to the
// checks.
func NewCheckSuite(secrets *Secrets, metadataClient voucher.MetadataClient, repositoryClient repository.Client, names ...string) (*voucher.Suite, error) {
	auth := newAuth()
	repos := validRepos()
	scanner := newScanner(secrets, metadataClient, auth)
	checksuite := voucher.NewSuite()

	trustedBuildCreators := viper.GetStringSlice("trusted_builder_identities")
	trustedProjects := viper.GetStringSlice("trusted_projects")

	checks, err := voucher.GetCheckFactories(names...)
	if nil != err {
		return checksuite, fmt.Errorf("can't create check suite: %s", err)
	}

	for name, check := range checks {
		setCheckAuth(check, auth)
		setCheckScanner(check, scanner)
		setCheckMetadataClient(check, metadataClient)
		setCheckValidRepos(check, repos)
		setCheckTrustedIdentitiesAndProjects(check, trustedBuildCreators, trustedProjects)
		setCheckRepositoryClient(check, repositoryClient)

		checksuite.Add(name, check)
	}

	return checksuite, nil
}
