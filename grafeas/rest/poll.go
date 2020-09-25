package rest

import (
	"context"
	"time"

	"github.com/Shopify/voucher"
	vgrafeas "github.com/Shopify/voucher/grafeas"
	"github.com/Shopify/voucher/grafeas/rest/objects"
)

const (
	attempts = 5
	sleep    = time.Second * 10
)

// isDone returns true if the passed discovery has finished, false otherwise.
func isDone(occurrence *objects.Occurrence) bool {
	occDiscovery := occurrence.Discovered
	if nil != occDiscovery {
		discovered := occDiscovery.Discovered
		if nil != discovered {
			if objects.DiscoveredAnalysisStatusFinishedSuccess == *discovered.AnalysisStatus {
				return true
			}
		}
	}

	return false
}

// pollForDiscoveries pauses execution until grafeas has pushed
// the Vulnerability information to the server.
func pollForDiscoveries(ctx context.Context, c *Client) error {
	for i := 0; i < attempts; i++ {
		discoveries, err := getVulnerabilityDiscoveries(ctx, c)
		if err != nil && !voucher.IsNoMetadataError(err) {
			return err
		}
		if len(discoveries) > 0 {
			for _, discoveryItem := range discoveries {
				if isDone(&discoveryItem) {
					return nil
				}
			}
		}
		time.Sleep(sleep)
	}
	return vgrafeas.ErrDiscoveriesUnfinished
}

func getVulnerabilityDiscoveries(ctx context.Context, g *Client) (items []objects.Occurrence, err error) {
	occurrences, err := g.getAllOccurrences(ctx)

	for _, occ := range occurrences {
		if *occ.Kind == objects.NoteKindDiscovery {
			items = append(items, occ)
		}
	}

	if 0 == len(items) && nil == err {
		err = &voucher.NoMetadataError{
			Type: vgrafeas.DiscoveryType,
			Err:  vgrafeas.ErrNoOccurrences,
		}
	}
	return
}
