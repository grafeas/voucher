package clair

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/voucher"
)

// Scanner implements the interface SnakeoilScanner.
type Scanner struct {
	config Config
	failOn voucher.Severity
	auth   voucher.Auth
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

	rawVulns, err := getVulnerabilities(ctx, scanner.config, tokenSrc, i)
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

// SetBasicAuth sets the username and password to use for Basic Auth,
// and enforces the use of Basic Auth for new connections.
func (scanner *Scanner) SetBasicAuth(username, password string) {
	scanner.config.Username = username
	scanner.config.Password = password
}

// NewScanner creates a new Scanner.
func NewScanner(config Config, auth voucher.Auth) *Scanner {
	scanner := new(Scanner)

	scanner.config = config

	scanner.auth = auth

	return scanner
}
