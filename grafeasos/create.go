package grafeasos

import (
	"github.com/docker/distribution/reference"
	grafeaspb "github.com/grafeas/client-go/0.1.0"
)

func (g *Client) getCreateOccurrence(reference reference.Canonical, parentNoteID string, attestation *grafeaspb.V1beta1attestationDetails, binauthProjectPath string) grafeaspb.V1beta1Occurrence {
	noteName := binauthProjectPath + "/notes/" + parentNoteID

	resource := grafeaspb.V1beta1Resource{
		Uri: "https://" + reference.Name() + "@" + reference.Digest().String(),
	}

	noteKind := grafeaspb.ATTESTATION_V1beta1NoteKind

	occurrence := grafeaspb.V1beta1Occurrence{Resource: &resource, NoteName: noteName, Kind: &noteKind, Attestation: attestation}

	return occurrence
}
