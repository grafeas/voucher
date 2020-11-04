package config

import (
	"strings"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/checks/org"
)

func RegisterDynamicChecks() {
	orgs := GetOrganizationsFromConfig()
	for alias, organization := range orgs {
		orgCheck := org.NewOrganizationCheckFactory(organization)
		voucher.RegisterCheckFactory("is_"+strings.ToLower(alias), orgCheck)
	}
}
