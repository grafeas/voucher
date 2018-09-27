package main

import (
	"github.com/Shopify/voucher/cmd/config"
	"github.com/Shopify/voucher/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the server",
	Long: `Run the go server on the specified port
	use --port=<port> to specify the port you want the server to run on`,
	Run: func(cmd *cobra.Command, args []string) {
		config := server.Config{
			Port:        viper.GetInt("server.port"),
			Timeout:     viper.GetInt("server.timeout"),
			RequireAuth: viper.GetBool("server.require_auth"),
			Username:    viper.GetString("server.username"),
			PassHash:    viper.GetString("server.password"),
		}
		server.Serve(&config)
	},
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	serverCmd.Flags().IntP("port", "", 8000, "port on which the server will listen")
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	serverCmd.Flags().StringVarP(&config.FileName, "config", "c", "", "path to config")
	serverCmd.Flags().IntP("timeout", "", 240, "number of seconds that should be dedicated to a Voucher call")
	viper.BindPFlag("server.timeout", serverCmd.Flags().Lookup("timeout"))
}
