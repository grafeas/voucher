package voucher

// AuthorizedCheck represents a Voucher check that needs to be authorized.
// For example, a check that needs to connect to the registry will
// need to implement AuthorizedCheck.
type AuthorizedCheck interface {
	Check
	SetAuth(Auth)
}
