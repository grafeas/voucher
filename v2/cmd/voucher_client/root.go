package main

import (
	"errors"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verify  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "voucher_client",
	Short: "voucher_client sends images to a Voucher server to be reviewed",
	Long: `voucher_client is a frontend for Voucher server, which allows users to send 
images for analysis. It automatically resolves tags to digests when it encounters
them.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing the image to check")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if verify {
			LookupAndVerify(args)
			return
		}
		LookupAndCheck(args)
	},
}

// init initializes the configuration and the flags.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolVar(&verify, "verify", false, "Verify instead of check an image.")
	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.voucher.yaml)")
	rootCmd.Flags().StringVarP(&defaultConfig.Server, "voucher", "v", "http://localhost:8000", "Voucher server to connect to.")
	viper.BindPFlag("server", rootCmd.Flags().Lookup("voucher"))
	rootCmd.Flags().StringVar(&defaultConfig.Username, "username", "", "Username to authenticate against Voucher with")
	viper.BindPFlag("username", rootCmd.Flags().Lookup("username"))
	rootCmd.Flags().StringVar(&defaultConfig.Password, "password", "", "Password to authenticate against Voucher with")
	viper.BindPFlag("password", rootCmd.Flags().Lookup("password"))
	rootCmd.Flags().IntVarP(&defaultConfig.Timeout, "timeout", "t", 240, "number of seconds to wait before failing")
	viper.BindPFlag("timeout", rootCmd.Flags().Lookup("timeout"))
	rootCmd.Flags().StringVarP(&defaultConfig.Check, "check", "c", "all", "the name of the checks to run against Voucher with")
	viper.BindPFlag("check", rootCmd.Flags().Lookup("check"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".vouch4cluster" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".voucher")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		if err = viper.Unmarshal(&defaultConfig); nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
