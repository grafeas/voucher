package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/cmd/config"
	"github.com/grafeas/voucher/repository"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (s *Server) handleChecks(w http.ResponseWriter, r *http.Request, name ...string) {
	var imageData voucher.ImageData
	var repositoryClient repository.Client
	var err error

	defer r.Body.Close()

	w.Header().Set("content-type", "application/json")

	LogRequests(r)

	imageData, err = handleInput(r)
	if nil != err {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		LogError(err.Error(), err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.serverConfig.TimeoutDuration())
	defer cancel()

	metadataClient, err := config.NewMetadataClient(ctx, s.secrets)
	if nil != err {
		http.Error(w, "server has been misconfigured", http.StatusInternalServerError)
		LogError("failed to create MetadataClient", err)
		return
	}
	defer metadataClient.Close()

	buildDetail, err := metadataClient.GetBuildDetail(ctx, imageData)
	if nil != err {
		LogWarning(fmt.Sprintf("could not get image metadata for %s", imageData), err)
	} else {
		if s.secrets != nil {
			repositoryClient, err = config.NewRepositoryClient(ctx, s.secrets.RepositoryAuthentication, buildDetail.RepositoryURL)
			if nil != err {
				LogWarning("failed to create repository client, continuing without git repo support:", err)
			}
		} else {
			log.Warning("failed to create repository client, no secrets configured")
		}
	}

	checksuite, err := config.NewCheckSuite(s.secrets, metadataClient, repositoryClient, name...)
	if nil != err {
		http.Error(w, "server has been misconfigured", http.StatusInternalServerError)
		LogError("failed to create CheckSuite", err)
		return
	}

	var results []voucher.CheckResult

	if viper.GetBool("dryrun") {
		results = checksuite.Run(ctx, s.metrics, imageData)
	} else {
		results = checksuite.RunAndAttest(ctx, metadataClient, s.metrics, imageData)
	}

	checkResponse := voucher.NewResponse(imageData, results)

	LogResult(checkResponse)

	err = json.NewEncoder(w).Encode(checkResponse)
	if nil != err {
		// if all else fails
		http.Error(w, err.Error(), http.StatusInternalServerError)
		LogError("failed to encode respoonse as JSON", err)
		return
	}
}
