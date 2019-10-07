package voucher

// ProvenanceCheck represents a Voucher check that sets
// trusted projects and build creators
type ProvenanceCheck interface {
	Check
	SetTrustedBuildCreators([]string)
	SetTrustedProjects([]string)
}
