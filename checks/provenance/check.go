package provenance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/repository"
)

// ErrNoBuildData is an error returned if we can't pull any BuildData from
// Grafeas for an image.
var ErrNoBuildData = errors.New("no build metadata associated with this image")

// check holds the required data for the check
type check struct {
	metadataClient       voucher.MetadataClient
	trustedBuildCreators map[string]bool
	trustedProjects      map[string]bool
}

// SetMetadataClient sets the MetadataClient for this Check.
func (p *check) SetMetadataClient(metadataClient voucher.MetadataClient) {
	p.metadataClient = metadataClient
}

// SetTrustedBuildCreators sets trustedBuildCreators for this Check.
func (p *check) SetTrustedBuildCreators(buildCreators []string) {
	p.trustedBuildCreators = make(map[string]bool)
	for _, identity := range buildCreators {
		p.trustedBuildCreators[identity] = true
	}
}

// SetTrustedProjects sets trustedProjects for this Check.
func (p *check) SetTrustedProjects(trustedProjects []string) {
	p.trustedProjects = make(map[string]bool)
	for _, project := range trustedProjects {
		p.trustedProjects[project] = true
	}
}

// Check runs the check :)
func (p *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	buildDetail, err := p.metadataClient.GetBuildDetail(ctx, i)
	if err != nil {
		if voucher.IsNoMetadataError(err) {
			return false, ErrNoBuildData
		}
		return false, err
	}

	if ok, err := validateProvenance(p, buildDetail); !ok || !validateArtifacts(i, buildDetail) {
		return false, err
	}

	return true, nil
}

func validateProvenance(p *check, detail repository.BuildDetail) (trusted bool, err error) {
	if !p.trustedBuildCreators[detail.BuildCreator] {
		err = fmt.Errorf("builder identity not trusted: %s", detail.BuildCreator)
		return
	}

	if !p.trustedProjects[detail.ProjectID] {
		err = fmt.Errorf("builder project not trusted: %s", detail.ProjectID)
		return
	}

	trusted = true
	return
}

func validateArtifacts(i voucher.ImageData, detail repository.BuildDetail) (matched bool) {
	// if an artifact built by this Build is the image, validate the SHAs match
	for _, artifact := range detail.Artifacts {
		if strings.HasSuffix(i.Digest().String(), artifact.Checksum) {
			matched = true
		}
	}
	return
}

func init() {
	voucher.RegisterCheckFactory("provenance", func() voucher.Check {
		return new(check)
	})
}
