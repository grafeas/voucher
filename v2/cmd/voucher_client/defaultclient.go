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
	switch strings.ToLower(defaultConfig.Auth) {
	case "idtoken":
		newClient, err := client.NewAuthClient(defaultConfig.Server)
		return newClient, err
	case "basic":
		newClient, err := client.NewClient(defaultConfig.Server)
		if err == nil {
			newClient.SetBasicAuth(defaultConfig.Username, defaultConfig.Password)
		}
		return newClient, err
	case "default-access-token":
		newClient, err := NewAuthClientWithToken(defaultConfig.Server)
		if err != nil {
			return nil, err
		}
		return newClient, err
	default:
		return nil, fmt.Errorf("invalid auth value: %q", defaultConfig.Auth)
	}
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(defaultConfig.Timeout)*time.Second)
}
