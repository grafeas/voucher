package containeranalysis

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/distribution/reference"
	"google.golang.org/api/iterator"

	grafeasv1 "cloud.google.com/go/grafeas/apiv1"
	grafeasapi "google.golang.org/genproto/googleapis/grafeas/v1"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/docker/uri"
	"github.com/Shopify/voucher/grafeas"
	vgrafeas "github.com/Shopify/voucher/grafeas"
)

const (
	caNoteProject       = "goog-analysis"
	caNoteID            = "PACKAGE_VULNERABILITY"
	discoPageSize int32 = 50
)

// vulnerabilityFilter returns a filter string
func vulnerabilityFilter(ref reference.Reference) string {
	return fmt.Sprintf(
		"%s AND noteProjectId=\"%s\" AND noteId=\"%s\"",
		vgrafeas.KindFilterStr(ref, grafeasapi.NoteKind_DISCOVERY),
		caNoteProject,
		caNoteID,
	)
}

func getVulnerabilityDiscoveries(ctx context.Context, client *grafeasv1.Client, ref reference.Reference) ([]*grafeasapi.DiscoveryOccurrence, error) {
	occurrences := make([]*grafeasapi.DiscoveryOccurrence, 0, 50)

	var err error

	project, err := uri.ReferenceToProjectName(ref)
	if nil != err {
		return nil, err
	}

	reqOccurrences := &grafeasapi.ListOccurrencesRequest{
		Parent:   vgrafeas.ProjectPath(project),
		Filter:   vulnerabilityFilter(ref),
		PageSize: discoPageSize,
	}

	occurrencesIterator := client.ListOccurrences(ctx, reqOccurrences)

	for {
		var occurrence *grafeasapi.Occurrence

		occurrence, err = occurrencesIterator.Next()
		if nil != err {
			if iterator.Done == err {
				err = nil
			}

			break
		}

		discoOccurrence := occurrence.GetDiscovery()
		if nil != discoOccurrence {
			occurrences = append(occurrences, discoOccurrence)
		}
	}

	if len(occurrences) == 0 && err == nil {
		err = &voucher.NoMetadataError{
			Type: grafeas.DiscoveryType,
			Err:  vgrafeas.ErrNoOccurrences,
		}
	}

	if err != nil {
		return nil, err
	}

	return occurrences, nil
}

// isDone returns true if the passed discovery has finished, false otherwise.
func isDone(occ *grafeasapi.DiscoveryOccurrence) bool {
	return occ.GetAnalysisStatus() == grafeasapi.DiscoveryOccurrence_FINISHED_SUCCESS
}

const attempts = 5
const sleep = time.Second * 10

// pollForDiscoveries pauses execution until Google Container Analysis has pushed
// the Vulnerability information to the server.
func pollForDiscoveries(ctx context.Context, c *Client, ref reference.Reference) error {
	for i := 0; i < attempts; i++ {
		discoveries, err := getVulnerabilityDiscoveries(
			ctx,
			c.containeranalysis,
			ref,
		)
		if err != nil && !voucher.IsNoMetadataError(err) {
			return fmt.Errorf("failed to get discoveries: %w", err)
		}

		if len(discoveries) > 0 {
			for _, discovery := range discoveries {
				if isDone(discovery) {
					return nil
				}
			}
		}

		time.Sleep(sleep)
	}

	return vgrafeas.ErrDiscoveriesUnfinished
}