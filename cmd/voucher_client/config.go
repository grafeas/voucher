package main

import (
	"context"
	"time"

	"github.com/Shopify/voucher/client"
)

type config struct {
	Hostname string
	Username string
	Password string
	Timeout  int
	Check    string
}

var defaultConfig = &config{}

func getCheck() string {
	return defaultConfig.Check
}

func getVoucherClient() (*client.VoucherClient, error) {
	newClient, err := client.NewClient(defaultConfig.Hostname, time.Duration(defaultConfig.Timeout)*time.Second)
	if nil == err {
		newClient.SetBasicAuth(defaultConfig.Username, defaultConfig.Password)
	}
	return newClient, err
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(defaultConfig.Timeout)*time.Second)
}
