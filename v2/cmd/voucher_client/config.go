package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/client"
	"google.golang.org/api/idtoken"
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
	options := []client.Option{
		client.WithUserAgent(fmt.Sprintf("voucher-client/%s", version)),
	}
	switch strings.ToLower(defaultConfig.Auth) {
	case "basic":
		options = append(options, client.WithBasicAuth(defaultConfig.Username, defaultConfig.Password))

	case "idtoken":
		idClient, err := idtoken.NewClient(context.Background(), defaultConfig.Server)
		if err != nil {
			return nil, err
		}
		options = append(options, client.WithHTTPClient(idClient))

	case "default-access-token":
		tokenClient, err := getDefaultTokenSourceClient(context.Background())
		if err != nil {
			return nil, err
		}
		options = append(options, client.WithHTTPClient(tokenClient))

	default:
		return nil, fmt.Errorf("invalid auth value: %q", defaultConfig.Auth)
	}

	newClient, err := client.NewClient(defaultConfig.Server, options...)
	if err != nil {
		return nil, err
	}
	return newClient, nil
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(defaultConfig.Timeout)*time.Second)
}
