package snakeoil

import (
	"context"
	"errors"

	"github.com/grafeas/voucher"
)

// ErrNoScanner is the error thrown when there is no SnakeoilScanner set for
// the Snakeoil Check.
var ErrNoScanner = errors.New("no scanner configured for snakeoil")

// check verifies if there are any known vulnerabilities for the
// passed image.
type check struct {
	scanner voucher.VulnerabilityScanner
}

// SetScanner sets the scanner that Snakeoil should use.
func (s *check) SetScanner(newScanner voucher.VulnerabilityScanner) {
	s.scanner = newScanner
}

// Check verifies if the image has known vulnerabilities
func (s *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
	if nil == s.scanner {
		return false, ErrNoScanner
	}

	vulns, err := s.scanner.Scan(ctx, i)
	if nil != err {
		return false, err
	}

	if 0 != len(vulns) {
		return false, voucher.NewVulnerabilityError(vulns)
	}

	return true, nil
}

func init() {
	voucher.RegisterCheckFactory("snakeoil", func() voucher.Check {
		return new(check)
	})
}
