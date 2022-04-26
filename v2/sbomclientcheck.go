package voucher

// SbomClientCheck represents a Voucher check that interacts
// with the sbom/gcr client.
type SbomClientCheck interface {
	Check
	SetSBOMClient(SBOMClient)
}
