package objects

import "time"

// NoteKind based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_note_kind.go
type NoteKind string

// consts
const (
	NoteKindUspecified    NoteKind = "NOTE_KIND_UNSPECIFIED"
	NoteKindVulnerability NoteKind = "VULNERABILITY"
	NoteKindBuild         NoteKind = "BUILD"
	NoteKindImage         NoteKind = "IMAGE"
	NoteKindPackage       NoteKind = "PACKAGE"
	NoteKindDeployment    NoteKind = "DEPLOYMENT"
	NoteKindDiscovery     NoteKind = "DISCOVERY"
	NoteKindAttestation   NoteKind = "ATTESTATION"
)

// Note based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_note.go
type Note struct {
	Name             string         `json:"name,omitempty"` //output only
	ShortDescription string         `json:"shortDescription,omitempty"`
	LongDescription  string         `json:"longDescription,omitempty"`
	Kind             *NoteKind      `json:"kind,omitempty"` //output only
	ExpirationTime   time.Time      `json:"expirationTime,omitempty"`
	CreateTime       time.Time      `json:"createTime,omitempty"` //output only
	UpdateTime       time.Time      `json:"updateTime,omitempty"` //output only
	RelatedNoteNames []string       `json:"relatedNoteNames,omitempty"`
	Vulnerability    *Vulnerability `json:"vulnerability,omitempty"`
	Build            *Build         `json:"build,omitempty"`
	Package          *Package       `json:"package,omitempty"`
	Discovery        *Discovery     `json:"discovery,omitempty"`
}
