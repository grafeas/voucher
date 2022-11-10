package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/go-containerregistry/pkg/v1/google"
	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/container"
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
		if err := clientRun(args); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}

func clientRun(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultConfig.Timeout)*time.Second)
	defer cancel()

	// Attempt google auth, but don't worry if it fails:
	googleAuth, _ := google.NewEnvAuthenticator()
	resolver := container.NewResolver(googleAuth)
	digest, err := resolver.ToDigest(ctx, args[0])
	if err != nil {
		return fmt.Errorf("getting canonical reference failed: %w", err)
	}

	client, err := getVoucherClient(ctx)
	if err != nil {
		return fmt.Errorf("creating client failed: %w", err)
	}
	var op func(context.Context, string, string) (voucher.Response, error)
	if verify {
		op = client.Verify
		fmt.Printf("Verifying %s\n", digest)
	} else {
		op = client.Check
		fmt.Printf("Checking %s\n", digest)
	}

	resp, err := op(ctx, getCheck(), digest.String())
	if err != nil {
		return fmt.Errorf("remote operation failed: %w", err)
	}
	fmt.Println(formatResponse(&resp))

	if !resp.Success {
		return fmt.Errorf("image failed to pass required check(s)")
	}

	return nil
}

// formatResponse returns the response as a string.
func formatResponse(resp *voucher.Response) string {
	output := ""
	if resp.Success {
		fmt.Println("image is approved")
	} else {
		fmt.Println("image was rejected")
	}
	for _, result := range resp.Results {
		if result.Success {
			output += fmt.Sprintf("   ✓ passed %s", result.Name)
			if !result.Attested {
				output += ", but wasn't attested"
			}
		} else {
			output += fmt.Sprintf("   ✗ failed %s", result.Name)
		}

		if "" != result.Err {
			output += fmt.Sprintf(", err: %s", result.Err)
		}
		output += "\n"
	}

	return output
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
	rootCmd.Flags().StringVarP(&defaultConfig.Auth, "auth", "a", "basic", "the method to authenticate against Voucher with. Supported types: basic, idtoken, default-access-token")
	viper.BindPFlag("auth", rootCmd.Flags().Lookup("auth"))
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
