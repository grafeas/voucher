package server

import (
	"net/http"

	"github.com/Shopify/voucher"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// LogRequests logs the request fields to stdout as Info
func LogRequests(r *http.Request) {
	r.ParseForm()
	log.WithFields(log.Fields{
		"url":  r.URL,
		"path": r.URL.Path,
		"form": r.Form,
	}).Info("Request Info")
}

// LogResult logs each test run as Info
func LogResult(response voucher.Response) {
	for _, result := range response.Results {
		log.WithFields(log.Fields{
			"check":    result.Name,
			"image":    response.Image,
			"passed":   result.Success,
			"attested": result.Attested,
			"error":    result.Err,
		}).Info("Check Result")
	}
}

// LogError logs server errors to stdout as Error
func LogError(err error) {
	log.Errorf("Server error: %s", err)
}
