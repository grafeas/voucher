package server

import (
	"encoding/json"
	"net/http"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/cmd/config"
	"github.com/spf13/viper"
)

func handle(w http.ResponseWriter, r *http.Request, name ...string) (err error) {

	var imageData voucher.ImageData

	defer r.Body.Close()

	w.Header().Set("content-type", "application/json")

	LogRequests(r)

	imageData, err = handleInput(r)
	if nil != err {
		http.Error(w, err.Error(), 422)
		LogError(err)
		return
	}

	metadataClient := config.NewMetadataClient()

	checksuite := config.NewCheckSuite(metadataClient, name...)

	var results []voucher.CheckResult

	if viper.GetBool("dryrun") {
		results = checksuite.Run(imageData)
	} else {
		results = checksuite.RunAndAttest(metadataClient, imageData)
	}

	checkResponse := voucher.NewResponse(imageData, results)

	LogResult(checkResponse)

	output, err := json.Marshal(checkResponse)
	if nil != err {
		// if all else fails
		http.Error(w, err.Error(), 500)
		LogError(err)
		return
	}
	w.Write(output)
	return
}

// HandleNobody is a request handler that makes the calls to create a "Nobody" attestation
func HandleNobody(w http.ResponseWriter, r *http.Request) { handle(w, r, "nobody") }

// HandleSnakeoil is a request handler that makes the calls to create a "Snakeoil" attestation
func HandleSnakeoil(w http.ResponseWriter, r *http.Request) { handle(w, r, "snakeoil") }

// HandleAll is a request handler that makes the calls to create all attestations, this includes DIY, Nobody, Snakeoil
func HandleAll(w http.ResponseWriter, r *http.Request) {
	handle(w, r, config.EnabledChecks(voucher.ToMapStringBool(viper.GetStringMap("checks")))...)
}

// HandleDIY is a request handler that makes the calls to create a "DIY" attestation
func HandleDIY(w http.ResponseWriter, r *http.Request) { handle(w, r, "diy") }

// HandleHealthCheck is a request handler that returns HTTP Status Code 200 when it is called from shopify cloud
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
