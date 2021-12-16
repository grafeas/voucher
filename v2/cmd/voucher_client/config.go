package main

import (
	"context"
	"fmt"
	"strings"
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
	var newClient *client.Client
	var err error
	switch strings.ToLower(defaultConfig.Auth) {
	case "idtoken":
		newClient, err = client.NewAuthClient(defaultConfig.Server)
		if err != nil {
			return nil, err
		}
	case "basic":
		newClient, err = client.NewClient(defaultConfig.Server)
		if err != nil {
			return nil, err
		}
		newClient.SetBasicAuth(defaultConfig.Username, defaultConfig.Password)
	case "default-access-token":
		newClient, err = NewAuthClientWithToken(context.Background(), defaultConfig.Server)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid auth value: %q", defaultConfig.Auth)
	}

	newClient.SetUserAgent(fmt.Sprintf("voucher-client/%s", version))
	return newClient, nil
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(defaultConfig.Timeout)*time.Second)
}
