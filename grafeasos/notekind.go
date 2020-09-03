package grafeasos

import (
	grafeaspb "github.com/grafeas/client-go/0.1.0"
)

type noteTest func(note *grafeaspb.V1beta1Note, noteKind grafeaspb.V1beta1NoteKind) bool

func isDiscoveryVulnerabilityNote(note *grafeaspb.V1beta1Note, noteKind grafeaspb.V1beta1NoteKind) bool {
	return *note.Kind == noteKind && *note.Discovery.AnalysisKind == grafeaspb.VULNERABILITY_V1beta1NoteKind
}

func isTypeNote(note *grafeaspb.V1beta1Note, noteKind grafeaspb.V1beta1NoteKind) bool {
	return *note.Kind == noteKind
}
