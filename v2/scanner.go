package voucher

import (
	"context"
)

// MetadataScanner implements voucher.VulnerabilityScanner, and connects to Grafeas
// to obtain vulnerability information.
type MetadataScanner struct {
	failOn Severity
	client MetadataClient
}

// FailOn sets severity level that a vulnerability must match or exheed to
// prompt a failure.
func (s *MetadataScanner) FailOn(severity Severity) {
	s.failOn = severity
}

// Scan gets the vulnerabilities for an Image.
func (s *MetadataScanner) Scan(ctx context.Context, i ImageData) ([]Vulnerability, error) {
	v, err := s.client.GetVulnerabilities(ctx, i)
	if nil != err {
		return []Vulnerability{}, err
	}
	vulns := make([]Vulnerability, 0, len(v))
	for _, item := range v {
		if ShouldIncludeVulnerability(item, s.failOn) {
			vulns = append(vulns, item)
		}
	}
	return vulns, nil
}

// NewScanner creates a new MetadataScanner.
func NewScanner(client MetadataClient) *MetadataScanner {
	scanner := new(MetadataScanner)
	scanner.client = client

	return scanner
}
