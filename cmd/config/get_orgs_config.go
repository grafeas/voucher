package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Shopify/voucher/repository"
)

func GetOrganizationsFromConfig() (orgs map[string]repository.Organization) {
	orgs = make(map[string]repository.Organization)
	repositories, ok := viper.Get("repositories").([]interface{})
	if !ok {
		repositories = []interface{}{}
	}
	for _, row := range repositories {
		if m, ok := row.(map[string]interface{}); ok {
			name := m["org-name"].(string)
			url := m["org-url"].(string)
			orgs[name] = repository.Organization{Name: name, URL: url}
		}
	}
	if len(orgs) == 0 {
		log.Warning("no repositories found")
	}
	return
}
