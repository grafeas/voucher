package voucher

import (
	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

// Occurrence is an alias for Google's containeranalysis Occurrence.
type Occurrence = *containeranalysispb.Occurrence

// NoteKind is an alias for Google's containeranalysis Note_Kind.
type NoteKind = containeranalysispb.Note_Kind
