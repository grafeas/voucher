package objects

import "time"

//NoteKind https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_note_kind.go
type NoteKind string

//consts
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

//Note https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_note.go
type Note struct {
	Name                 string                `json:"name,omitempty"` //output only
	ShortDescription     string                `json:"shortDescription,omitempty"`
	LongDescription      string                `json:"longDescription,omitempty"`
	Kind                 *NoteKind             `json:"kind,omitempty"` //output only
	RelatedURL           []RelatedURL          `json:"relatedUrl,omitempty"`
	ExpirationTime       time.Time             `json:"expirationTime,omitempty"`
	CreateTime           time.Time             `json:"createTime,omitempty"` //output only
	UpdateTime           time.Time             `json:"updateTime,omitempty"` //output only
	RelatedNoteNames     []string              `json:"relatedNoteNames,omitempty"`
	Vulnerability        *Vulnerability        `json:"vulnerability,omitempty"`
	Build                *Build                `json:"build,omitempty"`
	BaseImage            *ImageBasis           `json:"baseImage,omitempty"`
	Package              *Package              `json:"package,omitempty"`
	Deployable           *Deployable           `json:"deployable,omitempty"`
	Discovery            *Discovery            `json:"discovery,omitempty"`
	AttestationAuthority *AttestationAuthority `json:"attestationAuthority,omitempty"`
}

//Deployable https://github.com/grafeas/client-go/blob/master/0.1.0/model_deployment_deployable.go
type Deployable struct {
	ResourceURI []string `json:"resourceUri,omitempty"` //required
}

//ImageFingerprint https://github.com/grafeas/client-go/blob/master/0.1.0/model_image_fingerprint.go
type ImageFingerprint struct {
	V1Name string   `json:"v1Name,omitempty"` //required
	V2Blob []string `json:"v2Blob,omitempty"` //required
	V2Name string   `json:"v2Name,omitempty"` //output only
}

//ImageBasis https://github.com/grafeas/client-go/blob/master/0.1.0/model_image_basis.go
type ImageBasis struct {
	ResourceURL string            `json:"resourceUrl,omitempty"` //required
	Fingerprint *ImageFingerprint `json:"fingerprint,omitempty"` //required
}
