package objects

import (
	"time"

	"github.com/docker/distribution/reference"
)

//Occurrence based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_occurrence.go
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
	Discovered    *DiscoveryDetails     `json:"discovered,omitempty"`
	Attestation   *AttestationDetails   `json:"attestation,omitempty"`
}

//Resource based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_resource.go
type Resource struct {
	URI string `json:"uri,omitempty"` //required
}

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
