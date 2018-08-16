package main

import "github.com/spf13/cobra"

var nobodyCmd = &cobra.Command{
	Use:   "nobody",
	Short: "checks that container runs as non root",
	Long:  `Ensures that the container is running as a non-root user`,
	Run: func(cmd *cobra.Command, args []string) {
		runCheck("nobody")
	},
}

func init() {
	rootCmd.AddCommand(nobodyCmd)
}
