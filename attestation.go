package voucher

import (
	"github.com/grafeas/voucher/signer"
)

// Attestation is a structure that contains the Attestation data that we want
// to create an MetadataItem from.
type Attestation struct {
	CheckName string
	Body      string
}

// NewAttestation creates a new Attestation for the check with the passed name,
// with the payload as the body. The payload will then be signed by the key associated
// with the check (referenced by the checkName).
func NewAttestation(checkName string, payload string) Attestation {
	return Attestation{
		CheckName: checkName,
		Body:      payload,
	}
}

// SignedAttestation is a structure that contains the Attestation data as well
// as the signature and signing key ID.
type SignedAttestation struct {
	Attestation
	Signature string
	KeyID     string
}

// SignAttestation takes a keyring and attestation and signs the body of the
// payload with it, updating the Attestation's Signature field.
func SignAttestation(s signer.AttestationSigner, attestation Attestation) (SignedAttestation, error) {
	signature, keyID, err := s.Sign(attestation.CheckName, attestation.Body)
	if nil != err {
		return SignedAttestation{}, err
	}

	return SignedAttestation{
		Attestation: attestation,
		Signature:   signature,
		KeyID:       keyID,
	}, nil
}

// SignedAttestationToResult returns a CheckResults from the SignedAttestation
// passed to it. Check names is set as appropriate.
func SignedAttestationToResult(attestation SignedAttestation) CheckResult {
	return CheckResult{
		Name:     attestation.CheckName,
		Success:  true,
		Attested: true,
		Details:  attestation,
	}
}
