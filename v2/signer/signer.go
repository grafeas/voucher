package signer

type AttestationSigner interface {
	// Sign finds the key for a given check, signs the body and returns the signature and the key identifier
	Sign(checkName, body string) (string, string, error)
	Close() error
}
