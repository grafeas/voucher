package objects

import (
	"strings"
	"time"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/repository"
	"github.com/docker/distribution/reference"
)

//Occurrence adopted from https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_occurrence.go
type Occurrence struct {
	//output only, form: `projects/[PROJECT_ID]/occurrences/[OCCURRENCE_ID]
	Name          string                `json:"name,omitempty"`
	Resource      *Resource             `json:"resource,omitempty"` //required
	NoteName      string                `json:"noteName,omitempty"` //required, form: `projects/[PROVIDER_ID]/notes/[NOTE_ID]`
	Kind          *NoteKind             `json:"kind,omitempty"`     //output only
	Remediation   string                `json:"remediation,omitempty"`
	CreateTime    time.Time             `json:"createTime,omitempty"` //output only
	UpdateTime    time.Time             `json:"updateTime,omitempty"` //output only
	Vulnerability *VulnerabilityDetails `json:"vulnerability,omitempty"`
	Build         *BuildDetails         `json:"build,omitempty"`
	DerivedImage  *ImageDetails         `json:"derivedImage,omitempty"`
	Installation  *PackageDetails       `json:"installation,omitempty"`
	Deployment    *DeploymentDetails    `json:"deployment,omitempty"`
	Discovered    *DiscoveryDetails     `json:"discovered,omitempty"`
	Attestation   *AttestationDetails   `json:"attestation,omitempty"`
}

//Resource https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_resource.go
type Resource struct {
	URI string `json:"uri,omitempty"` //required
}

//RelatedURL https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_related_url.go
type RelatedURL struct {
	URL   string `json:"url,omitempty"`
	Label string `json:"label,omitempty"`
}

//not used, so not implemented types

//ImageDetails https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1image_details.go
type ImageDetails struct{}

//DeploymentDetails https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1deployment_details.go
type DeploymentDetails struct{}

//RPCStatus https://github.com/grafeas/client-go/blob/master/0.1.0/model_rpc_status.go
type RPCStatus struct{}

//NewOccurrence creates new occurrence
func NewOccurrence(reference reference.Canonical, parentNoteID string, attestation *AttestationDetails, binauthProjectPath string) Occurrence {
	noteName := binauthProjectPath + "/notes/" + parentNoteID

	resource := Resource{
		URI: "https://" + reference.Name() + "@" + reference.Digest().String(),
	}

	noteKind := NoteKindAttestation

	occurrence := Occurrence{Resource: &resource, NoteName: noteName, Kind: &noteKind, Attestation: attestation}

	return occurrence
}

//OccurrenceToAttestation converts objects.Occurrence to voucher.SignedAttestation
func OccurrenceToAttestation(checkName string, occ *Occurrence) voucher.SignedAttestation {
	signedAttestation := voucher.SignedAttestation{
		Attestation: voucher.Attestation{
			CheckName: checkName,
		},
	}

	attestationDetails := occ.Attestation

	signedAttestation.Body = string(*attestationDetails.Attestation.GenericSignedAttestation.ContentType)

	return signedAttestation
}

//OccurrenceToBuildDetail converts an Occurrence to a Build_Detail
func OccurrenceToBuildDetail(occ *Occurrence) (detail repository.BuildDetail) {
	buildProvenance := occ.Build.Provenance

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

// OccurrenceToVulnerability converts an Occurrence to a Vulnerability.
func OccurrenceToVulnerability(occ *Occurrence, vulProject string) (vul voucher.Vulnerability) {
	vul.Name = strings.Replace(occ.NoteName, vulProject, "", 1)
	vulnDetails := occ.Vulnerability

	vul.Severity = getSeverity(vulnDetails.Severity)
	packageIssues := vulnDetails.PackageIssue

	if vul.Severity == voucher.UnknownSeverity && len(packageIssues) > 0 {
		vul.Severity = getSeverity(vulnDetails.EffectiveSeverity)
	}

	return
}

// getSeverity translates the client-fo grafeas Severity to a Voucher Severity.
func getSeverity(severity *VulnerabilitySeverity) voucher.Severity {
	switch *severity {
	case SeverityMinimal:
		return voucher.NegligibleSeverity
	case SeverityLow:
		return voucher.LowSeverity
	case SeverityMedium:
		return voucher.MediumSeverity
	case SeverityHigh:
		return voucher.HighSeverity
	case SeverityCritical:
		return voucher.CriticalSeverity
	}
	return voucher.UnknownSeverity
}
