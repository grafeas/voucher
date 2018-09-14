package grafeas

import (
	"strings"
	"time"

	"github.com/Shopify/voucher"
)

// pollForDiscoveries pauses execution until Google Container Analysis has pushed
// the Vulnerability information to the server.
func pollForDiscoveries(c voucher.MetadataClient, i voucher.ImageData) error {
	attempts := 6
	for {
		if attempts < 0 {
			break
		}

		discoveries, err := c.GetMetadata(i, voucher.DiscoveryType)
		if err != nil {
			return err
		}

		if len(discoveries) > 0 {
			for _, discovery := range discoveries {
				noteName := strings.Split(discovery.NoteName, "/")
				noteKind := noteName[len(noteName)-1]

				if noteKind != "PACKAGE_VULNERABILITY" {
					continue
				}

				if discovery.GetDiscovered().GetOperation().GetDone() {
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
