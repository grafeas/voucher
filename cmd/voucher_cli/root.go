package main

import (
	"context"
	"time"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/cmd/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd stores the metadata about the toot voucher commands
var rootCmd = &cobra.Command{
	Use:   "voucher",
	Short: "An attestation generator for your containers",
	Long:  `voucher is a program that will create attestations for you Docker containers.`,
}

// Execute runs the cli commands
func Execute() {
	log.Println("exec")
	// setupConfig()
	rootCmd.Execute()
}

func getImageData() voucher.ImageData {
	image, err := rootCmd.PersistentFlags().GetString("image")
	if nil != err {
		log.Fatal(err)
	}

	imageData, err := voucher.NewImageData(image)
	if nil != err {
		log.Fatal(err)
	}
	return imageData
}

func runCheck(name ...string) {
	var results []voucher.CheckResult

	imageData := getImageData()

	context, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("timeout"))*time.Second)
	defer cancel()

	metadataClient := config.NewMetadataClient(context)

	checksuite, err := config.NewCheckSuite(metadataClient, name...)
	if nil != err {
		log.Fatalf("could not create CheckSuite: %s", err)
	}

	if viper.GetBool("dryrun") {
		results = checksuite.Run(imageData)
	} else {
		results = checksuite.RunAndAttest(metadataClient, imageData)
	}

	response := voucher.NewResponse(imageData, results)
	log.WithFields(log.Fields{
		"image":   response.Image,
		"success": response.Success,
	}).Info()

	for _, result := range response.Results {
		log.WithFields(log.Fields{
			"name":     result.Name,
			"success":  result.Success,
			"attested": result.Attested,
			"err":      result.Err,
		}).Info()
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	rootCmd.PersistentFlags().StringVarP(&config.FileName, "config", "c", "", "path to config")
	rootCmd.PersistentFlags().StringP("image", "i", "", "path of image to check")
	rootCmd.MarkPersistentFlagRequired("image")
	rootCmd.PersistentFlags().StringP("scanner", "", "", "vulnerability scanner to utilize")
	viper.BindPFlag("scanner", rootCmd.PersistentFlags().Lookup("scanner"))
	rootCmd.PersistentFlags().StringP("failon", "", "", "minimum vulnerability severity to fail on")
	viper.BindPFlag("failon", rootCmd.PersistentFlags().Lookup("failon"))
	rootCmd.PersistentFlags().BoolP("dryrun", "", false, "only run tests, do not push any attestations")
	viper.BindPFlag("dryrun", rootCmd.PersistentFlags().Lookup("dryrun"))
	rootCmd.PersistentFlags().IntP("timeout", "", 240, "number of seconds that should be dedicated to a Voucher call")
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	log.SetFormatter(&log.TextFormatter{})
}
