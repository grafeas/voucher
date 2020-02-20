package server

import (
	"net/http"

	"github.com/Shopify/voucher/cmd/config"
	"github.com/Shopify/voucher/metrics"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	serverConfig *Config
	secrets      *config.Secrets
	metrics      metrics.Client
}

// Serve creates a server on the specified port
func Serve(config *Config, secrets *config.Secrets, metrics metrics.Client) {
	s := &Server{config, secrets, metrics}
	router := NewRouter(s)
	log.Fatal(http.ListenAndServe(config.Address(), router))
}
