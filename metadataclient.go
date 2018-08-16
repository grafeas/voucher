package voucher

// MetadataClient is an interface that represents something that communicates
// with the Metadata server.
type MetadataClient interface {
	CanAttest() bool
	NewPayloadBody(ImageData) (string, error)
	GetOccurrencesForImage(ImageData, NoteKind) ([]Occurrence, error)
	AddAttestationOccurrenceToImage(ImageData, AttestationPayload) (Occurrence, error)
}
