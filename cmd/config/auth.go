package config

import (
	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/auth/google"
)

func newAuth() voucher.Auth {
	return google.NewAuth()
}
