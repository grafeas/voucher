package server

import (
	"net/http"

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
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range Routes {
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

// Routes an array of type Route
var Routes = []Route{
	{
		"All",
		"POST",
		"/all",
		HandleAll,
	},
	{
		"Individual Check",
		"POST",
		"/{check}",
		HandleIndividualCheck,
	},
	{
		"healthcheck: /services/ping",
		"GET",
		healthCheckPath,
		HandleHealthCheck,
	},
}
