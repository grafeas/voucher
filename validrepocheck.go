package voucher

// RepoValidatorCheck represents a Voucher check that validates the passed
// image is from a valid repo.
type RepoValidatorCheck interface {
	Check
	SetValidRepos(repos []string)
}
