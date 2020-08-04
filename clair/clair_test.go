package clair

import (
	"context"
	"net/http/httptest"
	"testing"

	v1 "github.com/coreos/clair/api/v1"
	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/Shopify/voucher"
	vtesting "github.com/Shopify/voucher/testing"
)

func TestGetClairVulnerabilitesSchema1(t *testing.T) {
	img, tokenSrc, clairServer := PrepareClairTest(t, ClairVulnerabilities())
	defer clairServer.Close()

	clairVulns, err := getClairVulnerabilities(vtesting.NewTestSchema1SignedManifest(vtesting.NewPrivateKey()), createClairConfig(clairServer.URL), tokenSrc, img)
	voucherVulns := convertToVoucherVulnerabilities(clairVulns, voucher.MediumSeverity)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(voucherVulns))
	require.ElementsMatch(t, VoucherVulnerabilities("medium", "high"), voucherVulns)
}

func TestGetClairVulnerabilitesSchema2(t *testing.T) {
	img, tokenSrc, clairServer := PrepareClairTest(t, ClairVulnerabilities())
	defer clairServer.Close()

	clairVulns, err := getClairVulnerabilities(vtesting.NewTestManifest(), createClairConfig(clairServer.URL), tokenSrc, img)
	voucherVulns := convertToVoucherVulnerabilities(clairVulns, voucher.MediumSeverity)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(voucherVulns))
	require.ElementsMatch(t, VoucherVulnerabilities("medium", "high"), voucherVulns)
}

func TestFilterClairVulnerabilities(t *testing.T) {
	img, tokenSrc, clairServer := PrepareClairTest(t, ClairVulnerabilities())
	defer clairServer.Close()

	clairVulns, err := getClairVulnerabilities(vtesting.NewTestManifest(), createClairConfig(clairServer.URL), tokenSrc, img)
	voucherVulns := convertToVoucherVulnerabilities(clairVulns, voucher.HighSeverity)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(voucherVulns))
	require.Equal(t, VoucherVulnerabilities("high"), voucherVulns)
}

func createClairConfig(hostname string) Config {
	return Config{
		Hostname: hostname,
		Username: "shopifolk",
		Password: "shopify",
	}
}

func PrepareClairTest(t *testing.T, clairVulns map[string][]v1.Vulnerability) (reference.Canonical, oauth2.TokenSource, *httptest.Server) {
	img := vtesting.NewTestReference(t)
	clairServer := vtesting.NewTestClairServer(t, clairVulns)
	tokenSrc, _ := vtesting.NewAuth(clairServer).GetTokenSource(context.TODO(), img)

	return img, tokenSrc, clairServer
}
