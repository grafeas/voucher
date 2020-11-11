package snakeoil

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafeas/voucher"
	vtesting "github.com/grafeas/voucher/testing"
)

func TestSnakeoilWithBadScanner(t *testing.T) {
	check := new(check)

	i, err := voucher.NewImageData("gcr.io/path/to/image@sha256:97db2bc359ccc94d3b2d6f5daa4173e9e91c513b0dcd961408adbb95ec5e5ce5")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)

	status, err := check.Check(context.Background(), i)

	require.Equalf(t, err, ErrNoScanner, "got wrong error for check: %s", err)
	assert.False(t, status, "check passed when it was not technically possible")
}

func TestSnakeoil(t *testing.T) {
	check := new(check)

	i, err := voucher.NewImageData("gcr.io/path/to/image@sha256:97db2bc359ccc94d3b2d6f5daa4173e9e91c513b0dcd961408adbb95ec5e5ce5")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)

	check.SetScanner(vtesting.NewScanner(t))

	status, err := check.Check(context.Background(), i)
	assert.NoErrorf(t, err, "check failed with error: %s", err)

	assert.True(t, status, "check failed when it should have passed")
}

func TestSnakeoilWithVulnerabilities(t *testing.T) {
	check := new(check)

	i, err := voucher.NewImageData("gcr.io/path/to/image@sha256:97db2bc359ccc94d3b2d6f5daa4173e9e91c513b0dcd961408adbb95ec5e5ce5")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)

	scanner := vtesting.NewScanner(t,
		voucher.Vulnerability{
			Name:        "cve-the-worst",
			Description: "it's really bad",
			Severity:    voucher.CriticalSeverity,
		},
		voucher.Vulnerability{
			Name:        "cve-this-is-fine",
			Description: "it's fine",
			Severity:    voucher.NegligibleSeverity,
		},
	)

	check.SetScanner(scanner)

	status, err := check.Check(context.Background(), i)
	require.Error(t, err, "check returned no errors, when it should have")

	assert.Truef(t, strings.HasPrefix(err.Error(), "vulnernable to 2 vulnerabilities:"), "error message is incorrectly formatted: %s", err)
	assert.Containsf(t, err.Error(), "cve-the-worst (critical)", "error message is incorrectly formatted: %s", err)
	assert.Containsf(t, err.Error(), "cve-this-is-fine (negligible)", "error message is incorrectly formatted: %s", err)
	assert.False(t, status, "check passed when it should have failed")
}
