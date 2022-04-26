package config

import (
	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/sbomgcr"
)

func newSBOMClient() voucher.SBOMClient {
	return sbomgcr.NewClient()
}
