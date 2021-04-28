package config

import (
	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/auth/google"
)

func newAuth() voucher.Auth {
	return google.NewAuth()
}
