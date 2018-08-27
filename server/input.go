package server

import (
	"encoding/json"
	"net/http"

	"github.com/Shopify/voucher"
)

// VoucherParams describes the input structure.
type VoucherParams struct {
	ImageURL string `json:"image_url"`
}

func handleInput(r *http.Request) (imageData voucher.ImageData, err error) {
	var params VoucherParams

	err = json.NewDecoder(r.Body).Decode(&params)
	if nil != err {
		return
	}

	imageData, err = voucher.NewImageData(params.ImageURL)
	return
}
