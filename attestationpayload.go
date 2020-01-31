package voucher

import (
	"github.com/Shopify/voucher/signer"
)

// AttestationPayload is a structure that contains the Attestation data that we
// want to create an MetadataItem from.
type AttestationPayload struct {
	CheckName string
	Body      string
}

// Sign takes a keyring and signs the body of the payload with it, returning that as a string.
func (payload *AttestationPayload) Sign(s signer.AttestationSigner) (string, string, error) {
	return s.Sign(payload.CheckName, payload.Body)
}

// NewAttestationPayload creates a new AttestationPayload for the check with the passed name,
// with the payload as the body. The payload will then be signed by the key associated
// with the check (referenced by the checkName).
func NewAttestationPayload(checkName string, payload string) AttestationPayload {
	return AttestationPayload{
		CheckName: checkName,
		Body:      payload,
	}
}
