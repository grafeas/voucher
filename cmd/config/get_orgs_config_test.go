package config

import (
	"testing"

	"github.com/Shopify/voucher/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetOrganizationsFromConfig(t *testing.T) {
	FileName = "../../testdata/config.toml"
	InitConfig()

	expected := map[string]repository.Organization{
		"shopify": repository.Organization{
			Alias: "shopify",
			Name:  "Shopify",
			VCS:   "github.com",
		},
	}
	got := GetOrganizationsFromConfig()

	assert.Equal(t, expected, got)
}
