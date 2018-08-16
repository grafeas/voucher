package main

import (
	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/cmd/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Runs all checks",
	Long:  `Runs all checks. Usage: voucher all -i your-image-here -p project`,
	Run: func(cmd *cobra.Command, args []string) {
		runCheck(config.EnabledChecks(voucher.ToMapStringBool(viper.GetStringMap("checks")))...)
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}
