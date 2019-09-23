package grafeas

import (
	"context"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/discovery"

	"github.com/Shopify/voucher"
)

// isVulnerabilityDiscovery returns true if the passed Item wraps a Vulnerability Discovery.
func isVulnerabilityDiscovery(item *Item) bool {
	return strings.Contains(item.Occurrence.NoteName, "VULNERABILITY")
}

// isDone returns true if the passed discovery has finished, false otherwise.
func isDone(item *Item) bool {
	occDiscovery := item.Occurrence.GetDiscovered()
	if nil != occDiscovery {
		discovered := occDiscovery.GetDiscovered()
		if nil != discovered {
			if discovery.Discovered_FINISHED_SUCCESS == discovered.GetAnalysisStatus() {
				return true
			}
		}
	}

	return false
}

const attempts = 5
const sleep = time.Second * 10

// pollForDiscoveries pauses execution until Google Container Analysis has pushed
// the Vulnerability information to the server.
func pollForDiscoveries(ctx context.Context, c voucher.MetadataClient, img voucher.ImageData) error {
	for i := 0; i < attempts; i++ {
		discoveries, err := c.GetMetadata(ctx, img, DiscoveryType)
		if err != nil {
			return err
		}
		if len(discoveries) > 0 {
			for _, discoveryItem := range discoveries {
				item, ok := discoveryItem.(*Item)
				if !ok || !isVulnerabilityDiscovery(item) {
					continue
				}
				if isDone(item) {
					return nil
				}
			}
		}
		time.Sleep(sleep)
	}
	return errDiscoveriesUnfinished
}
