package vtesting

import (
	"context"
	"testing"

	"github.com/grafeas/voucher"
)

type testVulnerabilityScanner struct {
	vulnerabilities []voucher.Vulnerability
}

// setVulnerabilitites sets the vulnerabilities to fail with.
func (t *testVulnerabilityScanner) setVulnerabilities(vulnerabilities []voucher.Vulnerability) {
	t.vulnerabilities = vulnerabilities
}

// FailOn sets the minimum Severity to consider an image vulnerable.
func (t *testVulnerabilityScanner) FailOn(failOn voucher.Severity) {
	// noop because the passed vulnerabilities will be the ones that are returned.
}

// Scan runs a scan against the passed ImageData and returns a slice of
// Vulnerabilities.
func (t *testVulnerabilityScanner) Scan(ctx context.Context, i voucher.ImageData) ([]voucher.Vulnerability, error) {
	return t.vulnerabilities, nil
}

// NewScanner creates a new Scanner suitable for testing with. The scanner will return
// all of the vulnerabilities that were passed in, regardless of what FailOn is set to. If
// the scanner is created with 0 vulnerabilities, Checks that use it will always pass.
func NewScanner(t *testing.T, vulnerabilities ...voucher.Vulnerability) voucher.VulnerabilityScanner {
	scanner := new(testVulnerabilityScanner)
	scanner.setVulnerabilities(vulnerabilities)
	return scanner
}
