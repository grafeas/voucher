package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Shopify/voucher"
)

// inputParams describes the input structure.
type inputParams struct {
	ImageURL string `json:"image_url"`
}

func handleInput(r *http.Request) (imageData voucher.ImageData, err error) {
	var body []byte

	// Read body
	body, err = ioutil.ReadAll(r.Body)
	if nil != err {
		return
	}

	var params inputParams

	// Unmarshal
	err = json.Unmarshal(body, &params)
	if nil != err {
		return
	}

	imageData, err = voucher.NewImageData(params.ImageURL)
	return
}
