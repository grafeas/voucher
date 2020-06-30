package clair

import (
	"context"
	"errors"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"golang.org/x/oauth2"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/docker"
	"github.com/Shopify/voucher/docker/schema2"
)

var errUnsupportedManifest = errors.New("manifests that are not schema version 2 are unsupported")

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
func (scanner *Scanner) Scan(ctx context.Context, i voucher.ImageData) ([]voucher.Vulnerability, error) {
	vulns := make([]voucher.Vulnerability, 0)

	tokenSrc, err := scanner.auth.GetTokenSource(ctx, i)
	if nil != err {
		return vulns, err
	}

	manifest, err := getDockerManifest(ctx, tokenSrc, i)
	if nil != err {
		return vulns, err
	}

	if !schema2.IsManifest(manifest) {
		return vulns, errUnsupportedManifest
	}

	s2Manifest := schema2.ToManifest(manifest)

	clairVulns, err := getClairVulnerabilities(s2Manifest, scanner.config, tokenSrc, i)
	if nil != err {
		return vulns, err
	}

	return convertToVoucherVulnerabilities(clairVulns, scanner.failOn), nil
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

func getDockerManifest(ctx context.Context, tokenSrc oauth2.TokenSource, image reference.Canonical) (distribution.Manifest, error) {
	client := oauth2.NewClient(ctx, tokenSrc)
	return docker.RequestManifest(client, image)
}
