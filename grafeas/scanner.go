package grafeas

import (
	"context"

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

// Scan gets the vulnerabilities for an Image.
func (s *Scanner) Scan(ctx context.Context, i voucher.ImageData) ([]voucher.Vulnerability, error) {
	v, err := s.client.GetVulnerabilities(ctx, i)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}
	vulns := make([]voucher.Vulnerability, 0, len(v))
	for _, item := range v {
		if voucher.ShouldIncludeVulnerability(item, s.failOn) {
			vulns = append(vulns, item)
		}
	}
	return vulns, nil
}

// NewScanner creates a new grafeas.Scanner.
func NewScanner(client voucher.MetadataClient) *Scanner {
	scanner := new(Scanner)
	scanner.client = client

	return scanner
}
