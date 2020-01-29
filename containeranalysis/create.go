package containeranalysis

import (
	grafeas "google.golang.org/genproto/googleapis/grafeas/v1"
)

func newOccurrenceAttestation(payload, signature, keyID string) *grafeas.Occurrence_Attestation {
	newAttestation := grafeas.AttestationOccurrence{
		SerializedPayload: []byte(payload),
		Signatures: []*grafeas.Signature{
			{
				Signature:   []byte(signature),
				PublicKeyId: keyID,
			},
		},
	}

	occurrenceAttestation := grafeas.Occurrence_Attestation{
		Attestation: &newAttestation,
	}

	return &occurrenceAttestation
}
