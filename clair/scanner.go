package clair

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/voucher"
)

// Scanner implements the interface SnakeoilScanner.
type Scanner struct {
	hostname string
	failOn   voucher.Severity
	auth     voucher.Auth
}

// FailOn sets severity level that a vulnerability must match or exheed to
// prompt a failure.
func (scanner *Scanner) FailOn(severity voucher.Severity) {
	scanner.failOn = severity
}

// Scan runs a scan in the Clair namespace.
func (scanner *Scanner) Scan(i voucher.ImageData) ([]voucher.Vulnerability, error) {
	vulns := make([]voucher.Vulnerability, 0)

	// We set a longer timeout for this, given that this operation is far more
	// intensive.
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	tokenSrc, err := scanner.auth.GetTokenSource(ctx, i)
	if nil != err {
		return vulns, err
	}

	rawVulns, err := getVulnerabilities(ctx, scanner.hostname, tokenSrc, i)
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
func NewScanner(hostname string, auth voucher.Auth) *Scanner {
	scanner := new(Scanner)

	scanner.hostname = hostname
	scanner.auth = auth

	return scanner
}
