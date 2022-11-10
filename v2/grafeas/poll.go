package grafeas

import (
	"context"
	"time"

	"github.com/docker/distribution/reference"
	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/container/gcr"
	"github.com/grafeas/voucher/v2/grafeas/objects"
)

var (
	attempts = 5
	sleep    = time.Second * 10
)

func setPollOptions(attemptsOption int, sleepOption time.Duration) {
	attempts = attemptsOption
	sleep = sleepOption
}

func defaultPollOptions() {
	attempts = 5
	sleep = time.Second * 10
}

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
func pollForDiscoveries(ctx context.Context, c *Client, ref reference.Reference) error {
	for i := 0; i < attempts; i++ {
		discoveries, err := getVulnerabilityDiscoveries(ctx, c, ref)
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

func getVulnerabilityDiscoveries(ctx context.Context, g *Client, ref reference.Reference) (items []objects.Occurrence, err error) {
	project, err := gcr.ReferenceToProjectName(ref)
	if nil != err {
		return nil, err
	}

	occurrences, err := g.getAllOccurrences(ctx, project)

	for _, occ := range occurrences {
		if *occ.Kind == objects.NoteKindDiscovery {
			items = append(items, occ)
		}
	}

	if 0 == len(items) && nil == err {
		err = &voucher.NoMetadataError{
			Type: DiscoveryType,
			Err:  errNoOccurrences,
		}
	}
	return
}
