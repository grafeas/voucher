package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	healthCheckPath     = "/services/ping"
	individualCheckPath = "/{check}"
	verifyCheckPath     = individualCheckPath + "/verify"
)

// Route stores metadata about a particular endpoint
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

// NewRouter creates a mux router with the specified routes and handlers
func NewRouter(s *Server) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range getRoutes(s) {
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func getRoutes(s *Server) []Route {
	return []Route{
		{
			"Check Image",
			"POST",
			individualCheckPath,
			s.HandleCheckImage,
		},
		{
			"Verify Image",
			"POST",
			verifyCheckPath,
			s.HandleVerifyImage,
		},
		{
			"healthcheck: /services/ping",
			"GET",
			healthCheckPath,
			s.HandleHealthCheck,
		},
	}
}
