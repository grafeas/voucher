package main

import "github.com/spf13/cobra"

var diyCmd = &cobra.Command{
	Use:   "diy",
	Short: "Ensures image was built by us",
	Long:  `Ensures image was build by us`,
	Run: func(cmd *cobra.Command, args []string) {
		runCheck("diy")
	},
}

func init() {
	rootCmd.AddCommand(diyCmd)
}
