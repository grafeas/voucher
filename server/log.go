package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/grafeas/voucher"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// LogRequests logs the request fields to stdout as Info
func LogRequests(r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.WithError(err).Info("received request with malformed form")
		return
	}

	log.WithFields(log.Fields{
		"url":  r.URL,
		"path": r.URL.Path,
		"form": r.Form,
	}).Info("received request")
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
func LogError(message string, err error) {
	log.Errorf("Server error: %s: %s", message, err)
}

// LogWarning logs server errors to stdout as Warning
func LogWarning(message string, err error) {
	log.Warningf("Server warning: %s: %s", message, err)
}

// LogInfo logs server information to stdout as Information.
func LogInfo(message string) {
	log.Infof("Server info: %s", message)
}
