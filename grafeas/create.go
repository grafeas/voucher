package grafeas

import (
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/attestation"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
)

func newOccurrenceAttestation(signature string, keyID string) *grafeaspb.Occurrence_Attestation {
	pgpKeyID := attestation.PgpSignedAttestation_PgpKeyId{
		PgpKeyId: keyID,
	}

	pgpSignedAttestation := attestation.PgpSignedAttestation{
		Signature:   signature,
		ContentType: attestation.PgpSignedAttestation_SIMPLE_SIGNING_JSON,
		KeyId:       &pgpKeyID,
	}

	attestationPgpSignedAttestation := attestation.Attestation_PgpSignedAttestation{
		PgpSignedAttestation: &pgpSignedAttestation,
	}

	newAttestation := attestation.Attestation{
		Signature: &attestationPgpSignedAttestation,
	}

	details := attestation.Details{
		Attestation: &newAttestation,
	}

	occurrenceAttestation := grafeaspb.Occurrence_Attestation{
		Attestation: &details,
	}

	return &occurrenceAttestation
}
