package objects

//AttestationSignedContentType https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_pgp_signed_attestation_content_type.go
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_generic_signed_attestation_content_type.go
type AttestationSignedContentType string

//consts
const (
	AttestationUnspecified AttestationSignedContentType = "CONTENT_TYPE_UNSPECIFIED"
	AttestationSigningJSON AttestationSignedContentType = "SIMPLE_SIGNING_JSON"
)

//AttestationDetails https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1attestation_details.go
type AttestationDetails struct {
	Attestation *Attestation `json:"attestation,omitempty"` //required
}

//Attestation https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_attestation.go
type Attestation struct {
	PgpSignedAttestation     *AttestationPgpSigned     `json:"pgpSignedAttestation,omitempty"`
	GenericSignedAttestation *AttestationGenericSigned `json:"genericSignedAttestation,omitempty"`
}

//AttestationPgpSigned https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_pgp_signed_attestation.go
type AttestationPgpSigned struct {
	ContentType *AttestationSignedContentType `json:"contentType,omitempty"`
	Signature   string                        `json:"signature,omitempty"` //required
	PgpKeyID    string                        `json:"pgpKeyId,omitempty"`
}

//AttestationGenericSigned https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_generic_signed_attestation.go
type AttestationGenericSigned struct {
	ContentType       *AttestationSignedContentType `json:"contentType,omitempty"`
	Signatures        []Signature                   `json:"signatures,omitempty"`
	SerializedPayload string                        `json:"serializedPayload,omitempty"`
}

//Signature https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_signature.go
type Signature struct {
	Signature   string `json:"signature,omitempty"`
	PublicKeyID string `json:"publicKeyId,omitempty"`
}

//attestation for note

//AttestationAuthority https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_authority.go
type AttestationAuthority struct {
	Hint *AuthorityHint `json:"hint,omitempty"`
}

//AuthorityHint https://github.com/grafeas/client-go/blob/master/0.1.0/model_authority_hint.go
type AuthorityHint struct {
	HumanReadableName string `json:"humanReadableName,omitempty"` //required
}
