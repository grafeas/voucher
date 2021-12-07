package main

import (
	"context"
	"time"

	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/client"
)

type config struct {
	Server   string
	Username string
	Password string
	Timeout  int
	Check    string
	Auth     string
}

var defaultConfig = &config{}

func getCheck() string {
	return defaultConfig.Check
}

func getVoucherClient() (voucher.Interface, error) {
	switch defaultConfig.Auth {
	case "idtoken":
		newClient, err := client.NewAuthClient(defaultConfig.Server)
		return newClient, err
	default:
		newClient, err := client.NewClient(defaultConfig.Server)
		if err == nil {
			newClient.SetBasicAuth(defaultConfig.Username, defaultConfig.Password)
		}
		return newClient, err
	}
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(defaultConfig.Timeout)*time.Second)
}
