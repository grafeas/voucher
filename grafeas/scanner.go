package grafeas

import (
	"fmt"

	"github.com/Shopify/voucher"
	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

// Scanner implements voucher.VulnerabilityScanner, and connects to Grafeas
// to obtain vulnerability information. It will block while scanning is active
// and fail if it spends more than a minute waiting for discovery to finish.
type Scanner struct {
	failOn voucher.Severity
	client voucher.MetadataClient
}

// FailOn sets severity level that a vulnerability must match or exheed to
// prompt a failure.
func (s *Scanner) FailOn(severity voucher.Severity) {
	s.failOn = severity
}

// GetVulnerabilitiesForImage returns the detected vulnerabilities for the Image
// described by voucher.ImageData.
func (s *Scanner) getVulnerabilitiesForImage(i voucher.ImageData) ([]voucher.Vulnerability, error) {
	occs, err := s.client.GetOccurrencesForImage(i, containeranalysispb.Note_PACKAGE_VULNERABILITY)
	vulns := make([]voucher.Vulnerability, 0, len(occs))
	if nil != err {
		return vulns, fmt.Errorf("could not get vulnerabilities: %s", err)
	}

	for _, occ := range occs {
		vuln := OccurrenceToVulnerability(occ)
		if voucher.ShouldIncludeVulnerability(vuln, s.failOn) {
			vulns = append(vulns, vuln)
		}
	}

	return vulns, nil
}

// Scan gets the vulnerabilities for an Image.
func (s *Scanner) Scan(i voucher.ImageData) ([]voucher.Vulnerability, error) {
	err := pollForDiscoveries(s.client, i)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}

	return s.getVulnerabilitiesForImage(i)
}

// NewScanner creates a new grafeas.Scanner.
func NewScanner(client voucher.MetadataClient) *Scanner {
	scanner := new(Scanner)
	scanner.client = client

	return scanner
}
