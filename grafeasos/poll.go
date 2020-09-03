package grafeasos

import (
	"context"
	"time"

	"github.com/Shopify/voucher"
	grafeaspb "github.com/grafeas/client-go/0.1.0"
)

// isDone returns true if the passed discovery has finished, false otherwise.
func isDone(occurrence *grafeaspb.V1beta1Occurrence) bool {
	occDiscovery := occurrence.Discovered
	if nil != occDiscovery {
		discovered := occDiscovery.Discovered
		if nil != discovered {
			if grafeaspb.FINISHED_SUCCESS_DiscoveredAnalysisStatus == *discovered.AnalysisStatus {
				return true
			}
		}
	}

	return false
}

const attempts = 5
const sleep = time.Second * 10

// pollForDiscoveries pauses execution until grafeas has pushed
// the Vulnerability information to the server.
func pollForDiscoveries(ctx context.Context, c *Client) error {
	for i := 0; i < attempts; i++ {
		discoveries, err := c.getVulnerabilityDiscoveries(ctx)
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
	return errDiscoveriesUnfinished
}
