package containeranalysis

import (
	"strings"

	grafeas "google.golang.org/genproto/googleapis/grafeas/v1"

	"github.com/Shopify/voucher"
	vgrafeas "github.com/Shopify/voucher/grafeas"
)

// OccurrenceToAttestation converts an Occurrence to a Attestation
func OccurrenceToAttestation(checkName string, occ *grafeas.Occurrence) voucher.SignedAttestation {
	signedAttestation := voucher.SignedAttestation{
		Attestation: voucher.Attestation{
			CheckName: checkName,
		},
	}

	attestationDetails := occ.GetAttestation()

	signedAttestation.Body = string(attestationDetails.GetSerializedPayload())

	return signedAttestation
}

func getCheckNameFromNoteName(project, value string) string {
	projectPath := vgrafeas.ProjectPath(project) + "/notes/"
	if strings.HasPrefix(value, projectPath) {
		result := strings.Replace(
			value,
			projectPath,
			"",
			-1,
		)
		if result != "" {
			return result
		}
	}
	return "unknown"
}
