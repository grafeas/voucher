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
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range getRoutes() {
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func getRoutes() []Route {
	return append(getCheckGroupRoutes(), Routes...)
}

// getCheckGroupRoutes creates Route objects for each group of required checks configured in the configuration file
func getCheckGroupRoutes() []Route {
	groups := config.GetRequiredChecksFromConfig()
	routes := make([]Route, 0, len(groups))
	for groupName := range groups {
		route := Route{
			Name:        groupName,
			Method:      "POST",
			Path:        "/" + groupName,
			HandlerFunc: HandleCheckGroup,
		}
		routes = append(routes, route)
	}
	return routes
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
