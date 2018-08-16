package clair

import (
	"fmt"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/docker"
)

// Scanner implements the interface SnakeoilScanner.
type Scanner struct {
	hostname string
	failOn   voucher.Severity
}

// FailOn sets severity level that a vulnerability must match or exheed to
// prompt a failure.
func (scanner *Scanner) FailOn(severity voucher.Severity) {
	scanner.failOn = severity
}

// Scan runs a scan in the Clair namespace.
func (scanner *Scanner) Scan(i voucher.ImageData) ([]voucher.Vulnerability, error) {
	vulns := make([]voucher.Vulnerability, 0)

	gcloudToken, err := voucher.GetAccessToken()
	if nil != err {
		return vulns, err
	}

	oauthToken, err := docker.Auth(gcloudToken, i)
	if nil != err {
		return vulns, err
	}

	rawVulns, err := getVulnerabilities(scanner.hostname, oauthToken, i)
	if nil != err {
		return vulns, err
	}

	vulns = make([]voucher.Vulnerability, 0, len(rawVulns))
	for _, rawVuln := range rawVulns {
		if "" == rawVuln.Name {
			fmt.Println("Empty Vulnerability???")
			continue
		}
		vuln := vulnerabilityToVoucherVulnerability(rawVuln)
		if voucher.ShouldIncludeVulnerability(vuln, scanner.failOn) {
			vulns = append(vulns, vuln)
		}
	}

	return vulns, err
}

// NewScanner creates a new Scanner.
func NewScanner(hostname string) *Scanner {
	scanner := new(Scanner)

	scanner.hostname = hostname

	return scanner
}
