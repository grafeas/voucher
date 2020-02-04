package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/cmd/config"
	"github.com/Shopify/voucher/repository"
)

func handleChecks(w http.ResponseWriter, r *http.Request, name ...string) {
	var imageData voucher.ImageData
	var repositoryClient repository.Client
	var err error

	defer r.Body.Close()

	if err = isAuthorized(r); nil != err {
		http.Error(w, "username or password is incorrect", 401)
		LogError("username or password is incorrect", err)
		return
	}

	w.Header().Set("content-type", "application/json")

	LogRequests(r)

	imageData, err = handleInput(r)
	if nil != err {
		http.Error(w, err.Error(), 422)
		LogError(err.Error(), err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverConfig.TimeoutDuration())
	defer cancel()

	metadataClient, err := config.NewMetadataClient(ctx)
	if nil != err {
		http.Error(w, "server has been misconfigured", 500)
		LogError("failed to create MetadataClient", err)
		return
	}
	defer metadataClient.Close()

	buildDetail, err := metadataClient.GetBuildDetail(ctx, imageData)
	if nil != err {
		LogWarning("failed to get buildDetail for image", err)
	} else {
		repositoryClient, err = config.NewRepositoryClient(ctx, buildDetail.RepositoryURL)
		if nil != err {
			LogWarning("failed to create repository client, continuing without git repo support:", err)
		}
	}

	checksuite, err := config.NewCheckSuite(metadataClient, repositoryClient, name...)
	if nil != err {
		http.Error(w, "server has been misconfigured", 500)
		LogError("failed to create CheckSuite", err)
		return
	}

	var results []voucher.CheckResult

	if viper.GetBool("dryrun") {
		results = checksuite.Run(ctx, imageData)
	} else {
		results = checksuite.RunAndAttest(ctx, metadataClient, imageData)
	}

	checkResponse := voucher.NewResponse(imageData, results)

	LogResult(checkResponse)

	err = json.NewEncoder(w).Encode(checkResponse)
	if nil != err {
		// if all else fails
		http.Error(w, err.Error(), 500)
		LogError("failed to encode respoonse as JSON", err)
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
