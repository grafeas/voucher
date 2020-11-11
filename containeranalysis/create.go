package containeranalysis

import (
	"github.com/docker/distribution/reference"
	grafeas "google.golang.org/genproto/googleapis/grafeas/v1"

	"github.com/grafeas/voucher"
)

func newOccurrenceAttestation(image reference.Canonical, attestation voucher.SignedAttestation, binauthProject string) *grafeas.CreateOccurrenceRequest {
	newAttestation := grafeas.AttestationOccurrence{
		SerializedPayload: []byte(attestation.Body),
		Signatures: []*grafeas.Signature{
			{
				Signature:   []byte(attestation.Signature),
				PublicKeyId: attestation.KeyID,
			},
		},
	}

	binauthProjectPath := projectPath(binauthProject)
	noteName := binauthProjectPath + "/notes/" + attestation.CheckName

	request := &grafeas.CreateOccurrenceRequest{
		Parent: binauthProjectPath,
		Occurrence: &grafeas.Occurrence{
			NoteName:    noteName,
			ResourceUri: "https://" + image.Name() + "@" + image.Digest().String(),
			Details: &grafeas.Occurrence_Attestation{
				Attestation: &newAttestation,
			},
		},
	}

	return request
}
