package config

import (
	"strings"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/checks/org"
)

func RegisterDynamicChecks() {
	orgs := GetOrganizationsFromConfig()
	for alias, organization := range orgs {
		orgCheck := org.NewOrganizationCheckFactory(organization)
		voucher.RegisterCheckFactory("is_"+strings.ToLower(alias), orgCheck)
	}
}
