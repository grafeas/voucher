package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/grafeas/voucher/cmd/config"
	"github.com/grafeas/voucher/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the server",
	Long: `Run the go server on the specified port
	use --port=<port> to specify the port you want the server to run on`,
	Run: func(cmd *cobra.Command, args []string) {
		serverConfig := server.Config{
			Port:        viper.GetInt("server.port"),
			Timeout:     viper.GetInt("server.timeout"),
			RequireAuth: viper.GetBool("server.require_auth"),
			Username:    viper.GetString("server.username"),
			PassHash:    viper.GetString("server.password"),
		}

		secrets, err := config.ReadSecrets()
		if err != nil {
			log.Printf("Error loading EJSON file, no secrets loaded: %v", err)
		}

		metricsClient, err := config.MetricsClient()
		if err != nil {
			log.Printf("Error configuring metrics client: %v", err)
		}

		config.RegisterDynamicChecks()

		voucherServer := server.NewServer(&serverConfig, secrets, metricsClient)

		for groupName, checks := range config.GetRequiredChecksFromConfig() {
			voucherServer.SetCheckGroup(groupName, checks)
		}

		voucherServer.Serve()
	},
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	serverCmd.Flags().IntP("port", "p", 8000, "port on which the server will listen")
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	serverCmd.Flags().StringVarP(&config.FileName, "config", "c", "", "path to config")
	serverCmd.Flags().IntP("timeout", "", 240, "number of seconds that should be dedicated to a Voucher call")
	viper.BindPFlag("server.timeout", serverCmd.Flags().Lookup("timeout"))
}
