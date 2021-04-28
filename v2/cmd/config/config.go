package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// FileName is the filename for the voucher configuration
var FileName string

// InitConfig searches for and loads configuration file for voucher
func InitConfig() {
	log.Println("initconfig")

	if FileName != "" {
		viper.SetConfigFile(FileName)
	} else {
		viper.SetConfigName("config")        // name of config file (without extension)
		viper.AddConfigPath("/etc/voucher/") // path to look for the config file in
		viper.AddConfigPath("$HOME/.voucher")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".") // optionally look for config in the working directory
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("config file: %s \n", err)
	}
	viper.AutomaticEnv()
}
