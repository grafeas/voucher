package grafeas

import (
	"fmt"

	"github.com/Shopify/voucher"
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
	items, err := s.client.GetMetadata(i, voucher.VulnerabilityType)
	vulns := make([]voucher.Vulnerability, 0, len(items))
	if nil != err {
		return vulns, fmt.Errorf("could not get vulnerabilities: %s", err)
	}

	for _, item := range items {
		metadataItem, ok := item.(*Item)
		if !ok {
			continue
		}

		vuln := OccurrenceToVulnerability(metadataItem.Occurrence)
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
