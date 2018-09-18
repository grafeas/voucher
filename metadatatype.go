package voucher

// MetadataType is a type which represents a MetadataClient's MetadataItem type.
type MetadataType string

const (
	// VulnerabilityType is specific to MetadataItem containing vulnerabilities.
	VulnerabilityType MetadataType = "vulnerability"
	// BuildDetailsType refers to MetadataItems containing image build details.
	BuildDetailsType MetadataType = "build details"
	// AttestationType refers to MetadataItems containing Binary Authorization Attestations.
	AttestationType MetadataType = "attestation"
)
