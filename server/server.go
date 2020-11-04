package server

import (
	"net/http"
	"strings"

	"github.com/grafeas/voucher/cmd/config"
	"github.com/grafeas/voucher/metrics"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	serverConfig *Config
	checkGroups  map[string][]string
	secrets      *config.Secrets
	metrics      metrics.Client
}

// NewServer creates a server on the specified port
func NewServer(config *Config, secrets *config.Secrets, metrics metrics.Client) *Server {
	return &Server{
		serverConfig: config,
		secrets:      secrets,
		metrics:      metrics,
		checkGroups:  make(map[string][]string),
	}
}

// Serve runs the Server on the specified port
func (server *Server) Serve() {
	router := NewRouter(server)
	log.Fatal(http.ListenAndServe(server.serverConfig.Address(), router))
}

// SetCheckGroup adds a list of checks as a group with the passed name.
func (server *Server) SetCheckGroup(name string, checkNames []string) {
	log.Infof("registering check group \"%s\": %s", name, strings.Join(checkNames, ", "))
	server.checkGroups[name] = checkNames
}

// HasCheckGroup returns true if the Check Group with the passed name has been
// registered with the server.
func (server *Server) HasCheckGroup(name string) bool {
	_, ok := server.checkGroups[name]
	return ok
}

// GetCheckGroup returns a list of checks names that are in the check group
// with the passed name.
func (server *Server) GetCheckGroup(name string) []string {
	checks := server.checkGroups[name]
	return checks
}
