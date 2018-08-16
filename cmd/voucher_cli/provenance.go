package main

import "github.com/spf13/cobra"

var provenanceCmd = &cobra.Command{
	Use:   "provenance",
	Short: "Ensures image was built by a trusted identity and that checksums match",
	Long:  `Ensures image was built by a trusted identity and that checksums match`,
	Run: func(cmd *cobra.Command, args []string) {
		runCheck("provenance")
	},
}

func init() {
	rootCmd.AddCommand(provenanceCmd)
}
