package voucher

// MetadataType is a type which represents a MetadataClient's Occurrence type.
type MetadataType string

const (
	VulnerabilityType MetadataType = "vulnerability"
	DiscoveryType     MetadataType = "discovery"
	BuildDetailsType  MetadataType = "build details"
)
