package config

import (
	"github.com/spf13/viper"

	"github.com/Shopify/voucher"
)

func GetRequiredChecksFromConfig() map[string][]string {
	requiredChecks := make(map[string][]string)
	requirements := viper.GetStringMap("required")
	if nil == requirements {
		requirements = map[string]interface{}{}
	}
	for env, val := range requirements {
		if m, ok := val.(map[string]interface{}); ok {
			requiredChecks[env] = EnabledChecks(voucher.ToMapStringBool(m))
		}
	}
	return requiredChecks
}
