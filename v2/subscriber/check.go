package subscriber

import (
	"context"

	"github.com/docker/distribution/reference"
	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/cmd/config"
	"github.com/grafeas/voucher/v2/repository"
)

// check runs all checks for a given image.
// Returns true if the required check(s) have passed and true if the check run needs to be retried.
func (s *Subscriber) check(canonicalImageReference reference.Canonical) (bool, bool) {
	var repositoryClient repository.Client
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.TimeoutDuration())
	defer cancel()

	metadataClient, err := config.NewMetadataClient(ctx, s.secrets)
	if nil != err {
		s.log.Errorf("failed to create MetadataClient: %s", err)
		return false, true
	}
	defer metadataClient.Close()

	buildDetail, err := metadataClient.GetBuildDetail(ctx, canonicalImageReference)
	if nil != err {
		s.log.Warningf("could not get image metadata for %s: %s", canonicalImageReference, err)
	} else {
		if s.secrets != nil {
			repositoryClient, err = config.NewRepositoryClient(ctx, s.secrets.RepositoryAuthentication, buildDetail.RepositoryURL)
			if nil != err {
				s.log.Warningf("failed to create repository client, continuing without git repo support: %s", err)
			}
		} else {
			s.log.Warning("failed to create repository client, no secrets configured")
		}
	}

	checksuite, err := config.NewCheckSuite(metadataClient, repositoryClient, s.cfg.RequiredChecks...)
	if nil != err {
		s.log.Errorf("failed to create CheckSuite: %s", err)
		return false, true
	}

	var results []voucher.CheckResult

	if s.cfg.DryRun {
		results = checksuite.Run(ctx, s.metrics, canonicalImageReference)
	} else {
		results = checksuite.RunAndAttest(ctx, metadataClient, s.metrics, canonicalImageReference)
	}

	checkResponse := voucher.NewResponse(canonicalImageReference, results)

	return checkResponse.Success, false
}
