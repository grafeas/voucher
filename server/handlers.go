package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/cmd/config"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func handleChecks(w http.ResponseWriter, r *http.Request, name ...string) {
	var imageData voucher.ImageData
	var err error

	defer r.Body.Close()

	if err = isAuthorized(r); nil != err {
		http.Error(w, "username or password is incorrect", 401)
		LogError(err)
		return
	}

	w.Header().Set("content-type", "application/json")

	LogRequests(r)

	imageData, err = handleInput(r)
	if nil != err {
		http.Error(w, err.Error(), 422)
		LogError(err)
		return
	}

	context, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	metadataClient := config.NewMetadataClient(context)

	checksuite, err := config.NewCheckSuite(metadataClient, name...)
	if nil != err {
		http.Error(w, "server has been misconfigured", 500)
		LogError(err)
		return
	}

	var results []voucher.CheckResult

	if viper.GetBool("dryrun") {
		results = checksuite.Run(imageData)
	} else {
		results = checksuite.RunAndAttest(metadataClient, imageData)
	}

	checkResponse := voucher.NewResponse(imageData, results)

	LogResult(checkResponse)

	err = json.NewEncoder(w).Encode(checkResponse)
	if nil != err {
		// if all else fails
		http.Error(w, err.Error(), 500)
		LogError(err)
		return
	}
}

// HandleAll is a request handler that makes the calls to create all attestations, this includes DIY, Nobody, Snakeoil
func HandleAll(w http.ResponseWriter, r *http.Request) {
	handleChecks(w, r, config.EnabledChecks(voucher.ToMapStringBool(viper.GetStringMap("checks")))...)
}

// HandleIndividualCheck is a request handler that executes an individual check and creates an attestation if applicable.
func HandleIndividualCheck(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)

	checkName := variables["check"]
	if "" == checkName {
		http.Error(w, "failure", 500)
		return
	}

	if voucher.IsCheckFactoryRegistered(checkName) {
		handleChecks(w, r, checkName)
		return
	}

	http.Error(w, fmt.Sprintf("check %s is not available", checkName), 404)
}

// HandleHealthCheck is a request handler that returns HTTP Status Code 200 when it is called from shopify cloud
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
