package server

import (
	"encoding/json"
	"net/http"

	voucher "github.com/grafeas/voucher/v2"
)

func handleInput(r *http.Request) (imageData voucher.ImageData, err error) {
	var request voucher.Request

	err = json.NewDecoder(r.Body).Decode(&request)
	if nil != err {
		return
	}

	imageData, err = voucher.NewImageData(request.ImageURL)
	return
}
