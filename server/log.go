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
	log.WithFields(log.Fields{
		"image":   response.Image,
		"results": response.Results,
	}).Info("Test Status")
}

// LogError logs server errors to stdout as Error
func LogError(err error) {
	log.Errorf("Server error: %s", err)
}
