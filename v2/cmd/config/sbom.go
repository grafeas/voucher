package config

import (
	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/sbomgcr"
)

func newSBOMClient() voucher.SBOMClient {
	service := sbomgcr.NewGCRService()
	return sbomgcr.NewClient(service)
}
