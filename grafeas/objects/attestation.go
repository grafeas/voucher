package objects

import (
	"github.com/grafeas/voucher"
)

//AttestationSignedContentType based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_pgp_signed_attestation_content_type.go
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_generic_signed_attestation_content_type.go
type AttestationSignedContentType string

//consts
const (
	AttestationUnspecified AttestationSignedContentType = "CONTENT_TYPE_UNSPECIFIED"
	AttestationSigningJSON AttestationSignedContentType = "SIMPLE_SIGNING_JSON"
)

//AttestationDetails based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1attestation_details.go
type AttestationDetails struct {
	Attestation *Attestation `json:"attestation,omitempty"` //required
}

//AsVoucherAttestation converts objects.AttestationDetails to voucher.SignedAttestation
func (ad *AttestationDetails) AsVoucherAttestation(checkName string) voucher.SignedAttestation {
	signedAttestation := voucher.SignedAttestation{
		Attestation: voucher.Attestation{
			CheckName: checkName,
		},
	}

	signedAttestation.Body = string(*ad.Attestation.GenericSignedAttestation.ContentType)

	return signedAttestation
}

//NewAttestation creates a new attestation
func NewAttestation(signedAttestation voucher.SignedAttestation) *AttestationDetails {
	contentType := AttestationSigningJSON
	return &AttestationDetails{Attestation: &Attestation{
		GenericSignedAttestation: &AttestationGenericSigned{
			Signatures: []Signature{{Signature: []byte(signedAttestation.Signature),
				PublicKeyID: signedAttestation.KeyID}}, ContentType: &contentType}}}
}

//Attestation based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_attestation.go
type Attestation struct {
	GenericSignedAttestation *AttestationGenericSigned `json:"genericSignedAttestation,omitempty"`
}

//AttestationGenericSigned based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_attestation_generic_signed_attestation.go
type AttestationGenericSigned struct {
	ContentType       *AttestationSignedContentType `json:"contentType,omitempty"`
	Signatures        []Signature                   `json:"signatures,omitempty"`
	SerializedPayload string                        `json:"serializedPayload,omitempty"`
}

//Signature based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_signature.go
type Signature struct {
	Signature   []byte `json:"signature,omitempty"`
	PublicKeyID string `json:"publicKeyId,omitempty"`
}
