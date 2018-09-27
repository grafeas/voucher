package grafeas

import (
	"strings"
	"time"

	"github.com/Shopify/voucher"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/discovery"
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

// pollForDiscoveries pauses execution until Google Container Analysis has pushed
// the Vulnerability information to the server.
func pollForDiscoveries(c voucher.MetadataClient, i voucher.ImageData) error {
	attempts := 6
	for {
		if attempts < 0 {
			break
		}

		discoveries, err := c.GetMetadata(i, DiscoveryType)
		if err != nil {
			return err
		}

		if len(discoveries) > 0 {
			for _, discoveryItem := range discoveries {

				item, ok := discoveryItem.(*Item)
				if !ok {
					continue
				}

				if !isVulnerabilityDiscovery(item) {
					continue
				}

				if isDone(item) {
					return nil
				}
				// TODO: add logging here.
				time.Sleep(time.Second * 10)
				attempts--
				break
			}
		} else {
			// TODO: add logging here as well.
			time.Sleep(time.Second * 5)
		}

	}
	return errDiscoveriesUnfinished
}
