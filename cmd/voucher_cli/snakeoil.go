package main

import "github.com/spf13/cobra"

var snakeoilCmd = &cobra.Command{
	Use:   "snakeoil",
	Short: "Ensures image passed vulnerability scan",
	Long:  `Ensures image passed vulnerability scan`,
	Run: func(cmd *cobra.Command, args []string) {
		runCheck("snakeoil")
	},
}

func init() {
	rootCmd.AddCommand(snakeoilCmd)
}
