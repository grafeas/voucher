package voucher

import (
	"context"
	"time"

	"github.com/grafeas/voucher/metrics"
)

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

// runner runs the passed check against the passed ImageData, and pushes results to the
// CheckResults channel.
func runner(ctx context.Context, name string, check Check, imageData ImageData, resultsChan chan CheckResult, metricsClient metrics.Client) {
	metricsClient.CheckRunStart(name)
	checkStart := time.Now()
	ok, err := check.Check(ctx, imageData)
	metricsClient.CheckRunLatency(name, time.Since(checkStart))
	if err == nil {
		if ok {
			metricsClient.CheckRunSuccess(name)
		} else {
			metricsClient.CheckRunFailure(name)
		}
		resultsChan <- CheckResult{Name: name, Err: "", Success: ok, ImageData: imageData}
	} else {
		metricsClient.CheckRunError(name, err)
		resultsChan <- CheckResult{Name: name, Err: err.Error(), Success: false, ImageData: imageData}
	}
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
func (cs *Suite) Run(ctx context.Context, metricsClient metrics.Client, imageData ImageData) []CheckResult {
	results := make([]CheckResult, 0, len(cs.checks))
	resultsChan := make(chan CheckResult, len(cs.checks))
	defer close(resultsChan)

	for name, check := range cs.checks {
		go runner(ctx, name, check, imageData, resultsChan, metricsClient)
	}

	for range cs.checks {
		results = append(results, <-resultsChan)
	}

	return results
}

// Attest runs through the passed []CheckResult and if a CheckResult is marked as successful,
// runs the CreateAttestion function in the Check corresponding to that CheckResult. Each
// CheckResult is updated with the details (or error) and the resulting []CheckResult is
// returned.
func (cs *Suite) Attest(ctx context.Context, metricsClient metrics.Client, metadataClient MetadataClient, results []CheckResult) []CheckResult {
	for i, result := range results {
		checkStart := time.Now()
		metricsClient.CheckAttestationStart(result.Name)
		if result.Success {
			details, err := createAttestation(ctx, metadataClient, result)
			results[i].Details = details
			if nil == err {
				results[i].Attested = true
				metricsClient.CheckAttestationSuccess(result.Name)
			} else {
				metricsClient.CheckAttestationError(result.Name, err)
				results[i].Err = err.Error()
			}
		}
		metricsClient.CheckAttestationLatency(result.Name, time.Since(checkStart))
	}

	return results
}

// RunAndAttest calls Run, followed by Attest, and returns the final []CheckResult.
func (cs *Suite) RunAndAttest(ctx context.Context, metadataClient MetadataClient, metricsClient metrics.Client, imageData ImageData) []CheckResult {
	results := cs.Run(ctx, metricsClient, imageData)
	return cs.Attest(ctx, metricsClient, metadataClient, results)
}

// createAttestation generates an attestation for the image Check described by CheckResult.
// That attestation is then added to the metadata server the MetadataClient is connected to.
func createAttestation(ctx context.Context, client MetadataClient, result CheckResult) (interface{}, error) {
	payload, err := client.NewPayloadBody(result.ImageData)
	if err != nil {
		return nil, err
	}

	attestation := NewAttestation(result.Name, payload)
	details, err := client.AddAttestationToImage(ctx, result.ImageData, attestation)
	return details, err
}

// NewSuite creates a new Suite.
func NewSuite() *Suite {
	suite := new(Suite)
	suite.checks = make(map[string]Check)
	return suite
}
