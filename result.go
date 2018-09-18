package voucher

// CheckResult describes the result of a Check. If a check failed, it will have a
// status of false. If a check succeeded, but its Attestation creation failed,
// Succes will be true, Attested will be false. Err will contain the first error to
// occur.
type CheckResult struct {
	ImageData ImageData    `json:"-"`
	Name      string       `json:"name"`
	Err       string       `json:"error,omitempty"`
	Success   bool         `json:"success"`
	Attested  bool         `json:"attested"`
	Details   MetadataItem `json:"details,omitempty"`
}
