package server

import (
	"encoding/json"
	"net/http"

	"github.com/Shopify/voucher"
)

// inputParams describes the input structure.
type inputParams struct {
	ImageURL string `json:"image_url"`
}

func handleInput(r *http.Request) (imageData voucher.ImageData, err error) {
	var params inputParams

	err = json.NewDecoder(r.Body).Decode(&params)
	if nil != err {
		return
	}

	imageData, err = voucher.NewImageData(params.ImageURL)
	return
}
