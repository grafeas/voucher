package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

var serverConfig *Config

// Serve creates a server on the specified port
func Serve(config *Config) {
	serverConfig = config
	router := NewRouter()
	log.Fatal(http.ListenAndServe(serverConfig.Address(), router))
}
