package grafeasos

import (
	"github.com/Shopify/voucher/repository"
	grafeaspb "github.com/grafeas/client-go/0.1.0"
)

// OccurrenceToBuildDetail converts an Occurrence to a Build_Detail
func OccurrenceToBuildDetail(occ *grafeaspb.V1beta1Occurrence) (detail repository.BuildDetail) {
	buildProvenance := occ.Build.Provenance

	detail.ProjectID = buildProvenance.ProjectId
	detail.BuildCreator = buildProvenance.Creator
	detail.BuildURL = buildProvenance.LogsUri
	detail.RepositoryURL = buildProvenance.SourceProvenance.Context.Git.Url
	detail.Commit = buildProvenance.SourceProvenance.Context.Git.RevisionId

	buildArtifacts := buildProvenance.BuiltArtifacts

	detail.Artifacts = make([]repository.BuildArtifact, 0, len(buildArtifacts))

	for _, artifact := range buildArtifacts {
		detail.Artifacts = append(detail.Artifacts, repository.BuildArtifact{
			ID:       artifact.Id,
			Checksum: artifact.Checksum,
		})
	}

	return
}
