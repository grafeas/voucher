package voucher

// MetadataCheck represents a Voucher check that interacts
// directly with a metadata server.
type MetadataCheck interface {
	Check
	SetMetadataClient(MetadataClient)
}
