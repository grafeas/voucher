package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Serve creates a server on the specified port
func Serve(port string) {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(port, router))
}
