package config

import (
	"github.com/spf13/viper"
)

func validRepos() []string {
	return viper.GetStringSlice("valid_repos")
}
