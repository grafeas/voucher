package voucher

// MetadataItem is a type which can be returned as a string.
type MetadataItem interface {
	String() string // String returns a string representation of the MetadataItem.
}
