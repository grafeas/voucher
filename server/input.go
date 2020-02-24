package server

import (
	"encoding/json"
	"net/http"

	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher"
)

func handleInput(r *http.Request) (imageData reference.Canonical, err error) {
	var request voucher.Request

	err = json.NewDecoder(r.Body).Decode(&request)
	if nil != err {
		return
	}

	imageData, err = voucher.NewImageReference(request.ImageURL)
	return
}
