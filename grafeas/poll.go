package grafeas

import (
	"time"

	"github.com/Shopify/voucher"
	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

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

				if item.Kind() != "PACKAGE_VULNERABILITY" {
					continue
				}

				if containeranalysispb.Discovery_Discovered_FINISHED_SUCCESS == item.Occurrence.GetDiscovered().AnalysisStatus {
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
