package voucher

// Suite is a suite of Checks, which
type Suite struct {
	checks map[string]Check
}

// Add adds a Check to the checks that can be run. Once a Check is added,
// it can be referenced by the name that was passed in when this function was called.
func (cs *Suite) Add(name string, check Check) {
	// Don't replace the existing check if it already exists.
	if !cs.Has(name) {
		cs.checks[name] = check
	}
}

// Has returns true if the passed check exists. Returns false if it does not.
func (cs *Suite) Has(name string) bool {
	return nil != cs.checks[name]
}

// Get returns the requested Check, or nil if one does not exist.
func (cs *Suite) Get(name string) (Check, error) {
	if cs.Has(name) {
		return cs.checks[name], nil
	}
	return nil, ErrNoCheck
}

// Run executes each of the Checks specified by the activeChecks parameter.
//
// For example, if a Suite has the "diy" and "nobody" tests, calling
//
//    Run(imageData)
//
// will run the "diy" and "nobody" tests.
//
// Run returns a []CheckResult with a CheckResult for each Check that was run.
func (cs *Suite) Run(imageData ImageData) []CheckResult {
	results := make([]CheckResult, 0, len(cs.checks))
	for name, check := range cs.checks {
		ok, err := check.Check(imageData)
		if err == nil {
			results = append(results, CheckResult{Name: name, Err: "", Success: ok, ImageData: imageData})
		} else {
			results = append(results, CheckResult{Name: name, Err: err.Error(), Success: false, ImageData: imageData})
		}
	}
	return results
}

// Attest runs through the passed []CheckResult and if a CheckResult is marked as successful,
// runs the CreateAttestion function in the Check corresponding to that CheckResult. Each
// CheckResult is updated with the details (or error) and the resulting []CheckResult is
// returned.
func (cs *Suite) Attest(metadataClient MetadataClient, results []CheckResult) []CheckResult {
	for i, result := range results {
		if result.Success {
			details, err := createAttestation(metadataClient, result)
			results[i].Details = details
			if nil == err {
				results[i].Attested = true
			} else {
				results[i].Err = err.Error()
			}
		}
	}
	return results
}

// RunAndAttest calls Run, followed by Attest, and returns the final []CheckResult.
func (cs *Suite) RunAndAttest(metadataClient MetadataClient, imageData ImageData) []CheckResult {
	results := cs.Run(imageData)
	return cs.Attest(metadataClient, results)
}

// createAttestation generates an attestation for the image Check described by CheckResult.
// That attestation is then added to the metadata server the MetadataClient is connected to.
func createAttestation(client MetadataClient, result CheckResult) (MetadataItem, error) {
	payload, err := client.NewPayloadBody(result.ImageData)
	if err != nil {
		return nil, err
	}

	attestationPayload := NewAttestationPayload(result.Name, payload)
	occ, err := client.AddAttestationToImage(result.ImageData, attestationPayload)
	return occ, err
}

// NewSuite creates a new Suite.
func NewSuite() *Suite {
	suite := new(Suite)
	suite.checks = make(map[string]Check)
	return suite
}
