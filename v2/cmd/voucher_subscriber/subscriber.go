package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/grafeas/voucher/v2/cmd/config"
	"github.com/grafeas/voucher/v2/subscriber"
)

var subscriberCmd = &cobra.Command{
	Use:   "subscriber",
	Short: "Runs the subscriber",
	Long: `Run the go subscriber that automatically vouches for images on the specified subscription and project
	use --project=<project> --subscription=<subscription> to specify the project and subscription you want the subscriber to pull from`,
	Run: func(cmd *cobra.Command, args []string) {
		var log = &logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		}

		secrets, err := config.ReadSecrets()
		if err != nil {
			log.Errorf("error loading EJSON file, no secrets loaded: %s", err)
		}

		metricsClient, err := config.MetricsClient()
		if err != nil {
			log.Errorf("error configuring metrics client: %s", err)
		}

		config.RegisterDynamicChecks()

		subscriberConfig := subscriber.Config{
			Project:        viper.GetString("pubsub.project"),
			Subscription:   viper.GetString("pubsub.subscription"),
			RequiredChecks: config.GetRequiredChecksFromConfig()["all"],
			DryRun:         viper.GetBool("dryrun"),
			Timeout:        viper.GetInt("pubsub.timeout"),
		}
		voucherSubscriber := subscriber.NewSubscriber(&subscriberConfig, secrets, metricsClient, log)

		err = voucherSubscriber.Subscribe(context.Background())
		if err != nil {
			log.Errorf("couldn't pull pub/sub messages: %s", err)
		}
	},
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	subscriberCmd.Flags().StringP("project", "p", "", "pub/sub project that has the subsciprion (required)")
	subscriberCmd.MarkFlagRequired("project")
	viper.BindPFlag("pubsub.project", subscriberCmd.Flags().Lookup("project"))
	subscriberCmd.Flags().StringP("subscription", "s", "", "pub/sub topic subscription (required)")
	subscriberCmd.MarkFlagRequired("subscription")
	viper.BindPFlag("pubsub.subscription", subscriberCmd.Flags().Lookup("subscription"))
	subscriberCmd.Flags().StringVarP(&config.FileName, "config", "c", "", "path to config")
	subscriberCmd.Flags().IntP("timeout", "", 240, "number of seconds that should be dedicated to a Voucher call")
	viper.BindPFlag("pubsub.timeout", subscriberCmd.Flags().Lookup("timeout"))
}
