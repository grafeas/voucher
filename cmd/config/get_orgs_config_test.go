package config

import (
	"testing"

	"github.com/grafeas/voucher/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetOrganizationsFromConfig(t *testing.T) {
	FileName = "../../testdata/config.toml"
	InitConfig()

	expected := map[string]repository.Organization{
		"shopify": {
			Alias: "shopify",
			Name:  "Shopify",
			VCS:   "github.com",
		},
	}
	got := GetOrganizationsFromConfig()

	assert.Equal(t, expected, got)
}
