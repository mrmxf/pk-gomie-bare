package walker

import (
	"llmProxy/gomie/logger"

	"github.com/gorilla/mux"
)

// The routes are walked in the order they were added. This function is called for each route. Sub-routers
// are explored depth-first.
func gorillaRouteWalker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	log := logger.GetLogger()
	path, _ := route.GetPathTemplate()
	log.Infow(path)
	return nil
}
