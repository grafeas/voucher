package server

import (
	"net/http"

	"github.com/Shopify/voucher/cmd/config"
	"github.com/gorilla/mux"
)

const healthCheckPath = "/services/ping"

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
	return append(getCheckGroupRoutes(s), []Route{
		{
			"All",
			"POST",
			"/all",
			s.HandleAll,
		},
		{
			"Individual Check",
			"POST",
			"/{check}",
			s.HandleIndividualCheck,
		},
		{
			"healthcheck: /services/ping",
			"GET",
			healthCheckPath,
			s.HandleHealthCheck,
		},
	}...,
	)
}

// getCheckGroupRoutes creates Route objects for each group of required checks configured in the configuration file
func getCheckGroupRoutes(s *Server) []Route {
	groups := config.GetRequiredChecksFromConfig()
	routes := make([]Route, 0, len(groups))
	for groupName := range groups {
		route := Route{
			Name:        groupName,
			Method:      "POST",
			Path:        "/" + groupName,
			HandlerFunc: s.HandleCheckGroup,
		}
		routes = append(routes, route)
	}
	return routes
}
