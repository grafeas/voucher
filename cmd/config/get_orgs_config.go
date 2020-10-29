package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/grafeas/voucher/repository"
)

func GetOrganizationsFromConfig() map[string]repository.Organization {
	orgs := make(map[string]repository.Organization)
	repositories := viper.GetStringMap("repository")
	if nil == repositories {
		repositories = map[string]interface{}{}
	}
	for alias, val := range repositories {
		if m, ok := val.(map[string]interface{}); ok {
			url := m["org-url"].(string)
			org := *repository.NewOrganization(alias, url)
			orgs[org.Alias] = *repository.NewOrganization(alias, url)
		}
	}
	if len(orgs) == 0 {
		log.Warning("no repositories found")
	}
	return orgs
}
