package grafeas

import (
	"github.com/Shopify/voucher"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
)

// OccurrenceToBuildDetails converts an Occurrence to a Build_Detail
func OccurrenceToBuildDetails(occ *grafeaspb.Occurrence) (detail voucher.BuildDetail) {
	buildProvenance := occ.GetBuild().GetProvenance()

	detail.ProjectID = buildProvenance.GetProjectId()
    detail.BuildCreator = buildProvenance.GetCreator()
	detail.BuildURL = buildProvenance.GetLogsUri()
	detail.RepositoryURL = buildProvenance.GetSourceProvenance().GetContext().GetGit().GetUrl()
	detail.Commit = buildProvenance.GetSourceProvenance().GetContext().GetGit().GetRevisionId()

	buildArtifacts := buildProvenance.GetBuiltArtifacts()

    detail.Artifacts = make([]voucher.BuildArtifact, 0, len(buildArtifacts))

    for _, artifact := range buildArtifacts {
        detail.Artifacts = append(detail.Artifacts, voucher.BuildArtifact{
            ID:       artifact.Id,
            Checksum: artifact.Checksum,
        })
    }

	return
}
