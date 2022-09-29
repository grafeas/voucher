package objects

import (
	"time"

	"github.com/grafeas/voucher/v2/repository"
)

//note objects

// Build based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_build_build.go
type Build struct {
	BuilderVersion string          `json:"builderVersion,omitempty"` //required
	Signature      *BuildSignature `json:"signature,omitempty"`
}

// BuildSignature based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_build_build_signature.go
type BuildSignature struct {
	PublicKey string `json:"publicKey,omitempty"`
	Signature string `json:"signature,omitempty"` //required
	KeyID     string `json:"keyId,omitempty"`
}

//occurrence objects

// BuildDetails based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1build_details.go
type BuildDetails struct {
	Provenance      *ProvenanceBuild `json:"provenance,omitempty"` //required
	ProvenanceBytes string           `json:"provenanceBytes,omitempty"`
}

// AsVoucherBuildDetail converts an BuildDetails to a Build_Detail
func (bd *BuildDetails) AsVoucherBuildDetail() (detail repository.BuildDetail) {
	buildProvenance := bd.Provenance

	detail.ProjectID = buildProvenance.ProjectID
	detail.BuildCreator = buildProvenance.Creator
	detail.BuildURL = buildProvenance.LogsURI
	detail.RepositoryURL = buildProvenance.SourceProvenance.Context.Git.URL
	detail.Commit = buildProvenance.SourceProvenance.Context.Git.RevisionID

	buildArtifacts := buildProvenance.BuiltArtifacts
	detail.Artifacts = make([]repository.BuildArtifact, 0, len(buildArtifacts))

	for _, artifact := range buildArtifacts {
		detail.Artifacts = append(detail.Artifacts, repository.BuildArtifact{
			ID:       artifact.ID,
			Checksum: artifact.Checksum,
		})
	}

	return
}

// ProvenanceBuild based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_build_provenance.go
type ProvenanceBuild struct {
	ID               string               `json:"id,omitempty"` //required
	ProjectID        string               `json:"projectId,omitempty"`
	BuiltArtifacts   []ProvenanceArtifact `json:"builtArtifacts,omitempty"`
	CreateTime       time.Time            `json:"createTime,omitempty"`
	StartTime        time.Time            `json:"startTime,omitempty"`
	EndTime          time.Time            `json:"endTime,omitempty"`
	Creator          string               `json:"creator,omitempty"` //email address
	LogsURI          string               `json:"logsUri,omitempty"`
	SourceProvenance *ProvenanceSource    `json:"sourceProvenance,omitempty"`
	TriggerID        string               `json:"triggerId,omitempty"`
	BuildOptions     map[string]string    `json:"buildOptions,omitempty"`
	BuilderVersion   string               `json:"builderVersion,omitempty"`
}

// ProvenanceArtifact based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_artifact.go
type ProvenanceArtifact struct {
	Checksum string   `json:"checksum,omitempty"`
	ID       string   `json:"id,omitempty"`
	Names    []string `json:"names,omitempty"`
}

// ProvenanceSource based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_source.go
type ProvenanceSource struct {
	ArtifactStorageSourceURI string          `json:"artifactStorageSourceUri,omitempty"`
	Context                  *SourceContext  `json:"context,omitempty"`
	AdditionalContexts       []SourceContext `json:"additionalContexts,omitempty"`
}

// SourceContext based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_source_context.go
type SourceContext struct {
	Git    *GitSourceContext `json:"git,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// GitSourceContext based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_git_source_context.go
type GitSourceContext struct { //SourceGitSourceContext
	URL        string `json:"url,omitempty"`
	RevisionID string `json:"revisionId,omitempty"`
}
