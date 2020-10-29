package config

import (
	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/auth/google"
)

func newAuth() voucher.Auth {
	return google.NewAuth()
}
