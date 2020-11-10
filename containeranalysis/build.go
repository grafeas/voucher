package containeranalysis

import (
	"github.com/Shopify/voucher/repository"
	grafeas "google.golang.org/genproto/googleapis/grafeas/v1"
)

// OccurrenceToBuildDetail converts an Occurrence to a BuildDetail
func OccurrenceToBuildDetail(occ *grafeas.Occurrence) (detail repository.BuildDetail) {
	buildProvenance := occ.GetBuild().GetProvenance()

	detail.ProjectID = buildProvenance.GetProjectId()
	detail.BuildCreator = buildProvenance.GetCreator()
	detail.BuildURL = buildProvenance.GetLogsUri()
	detail.RepositoryURL = buildProvenance.GetSourceProvenance().GetContext().GetGit().GetUrl()
	detail.Commit = buildProvenance.GetSourceProvenance().GetContext().GetGit().GetRevisionId()

	buildArtifacts := buildProvenance.GetBuiltArtifacts()

	detail.Artifacts = make([]repository.BuildArtifact, 0, len(buildArtifacts))

	for _, artifact := range buildArtifacts {
		detail.Artifacts = append(detail.Artifacts, repository.BuildArtifact{
			ID:       artifact.Id,
			Checksum: artifact.Checksum,
		})
	}

	return
}
