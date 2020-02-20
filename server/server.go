package server

import (
	"net/http"

	"github.com/Shopify/voucher/cmd/config"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	serverConfig *Config
	secrets      *config.Secrets
}

// Serve creates a server on the specified port
func Serve(config *Config, secrets *config.Secrets) {
	s := &Server{config, secrets}
	router := NewRouter(s)
	log.Fatal(http.ListenAndServe(config.Address(), router))
}
