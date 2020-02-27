package config

import (
	"github.com/spf13/viper"
)

func GetRequiredChecksFromConfig() map[string][]string {
	requiredChecks := make(map[string][]string)

	requiredChecks["all"] = toStringSlice(viper.GetStringMap("checks"))

	requirements := viper.GetStringMap("required")
	if nil == requirements {
		requirements = map[string]interface{}{}
	}
	for env, val := range requirements {
		if m, ok := val.(map[string]interface{}); ok {
			requiredChecks[env] = toStringSlice(m)
		}
	}
	return requiredChecks
}

// toStringSlice takes a map[string]interface{} and converts it to a
// slice of strings using the keys (dropping any values that do not cast to
// booleans cleanly, or have the value of false).
func toStringSlice(in map[string]interface{}) []string {
	out := make([]string, 0, len(in))
	for key, rawValue := range in {
		if value, ok := rawValue.(bool); ok {
			if value {
				out = append(out, key)
			}
		}
	}
	return out
}
