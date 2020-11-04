package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/cmd/config"
)

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request, names ...string) {
	var imageData voucher.ImageData
	var err error

	defer r.Body.Close()

	w.Header().Set("content-type", "application/json")

	LogRequests(r)

	imageData, err = handleInput(r)
	if nil != err {
		http.Error(w, err.Error(), 422)
		LogError(err.Error(), err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.serverConfig.TimeoutDuration())
	defer cancel()

	metadataClient, err := config.NewMetadataClient(ctx, s.secrets)
	if nil != err {
		http.Error(w, "server has been misconfigured", 500)
		LogError("failed to create MetadataClient", err)
		return
	}
	defer metadataClient.Close()

	attestations, err := metadataClient.GetAttestations(ctx, imageData)
	if nil != err {
		LogWarning(fmt.Sprintf("could not get image attestations for %s", imageData), err)
	}

	checkResponse := voucher.NewResponse(
		imageData,
		attestationsToResults(attestations, names),
	)

	LogResult(checkResponse)

	err = json.NewEncoder(w).Encode(checkResponse)
	if nil != err {
		// if all else fails
		http.Error(w, err.Error(), 500)
		LogError("failed to encode respoonse as JSON", err)
		return
	}
}

func attestationsToResults(attestations []voucher.SignedAttestation, names []string) []voucher.CheckResult {
	results := make([]voucher.CheckResult, 0, len(names))

	for _, name := range names {
		failed := true
		for _, attestation := range attestations {
			if attestation.CheckName == name {
				failed = false
				results = append(results, voucher.SignedAttestationToResult(attestation))
				break
			}
		}
		if failed {
			results = append(
				results,
				voucher.CheckResult{
					Name:     name,
					Err:      "",
					Success:  false,
					Attested: false,
					Details:  nil,
				},
			)
		}
	}

	return results
}
