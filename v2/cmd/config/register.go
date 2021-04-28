package config

import (
	"strings"

	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/checks/org"
)

func RegisterDynamicChecks() {
	orgs := GetOrganizationsFromConfig()
	for alias, organization := range orgs {
		orgCheck := org.NewOrganizationCheckFactory(organization)
		voucher.RegisterCheckFactory("is_"+strings.ToLower(alias), orgCheck)
	}
}
