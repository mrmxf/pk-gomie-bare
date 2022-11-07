package gomie

// Copyright (C) 2022 Dalet UK Ltd
// All rights reserved. Not for external use.

/**
 * @todo Add context.Context to handlers
 */
import (
	"llmProxy/gomie/config"
	"llmProxy/gomie/logger"
	"llmProxy/gomie/static"
	"llmProxy/limelm"

	"github.com/gorilla/mux"
	"github.com/mrmxf/pk-gomie-bare/config"
	"github.com/mrmxf/pk-gomie-bare/logger"
)

// Main high level controls for the proxyy server
type Proxy struct {
	Cfg    *config.Config // config object
	Router *mux.Router    // main router
	Log    *logger.Logger // pretty logger
}

// This function walks all the routes registered with the proxy server and then logs
// them to the console for debugging.
func (p *Proxy) gorillaWalkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, _ := route.GetPathTemplate()
	p.Log.Infow(p.Cfg.GetString("app_name")+": walkPaths ", "path", path)
	return nil
}

/* Main Entry Point
 */
func (p *Proxy) Setup(forceConfigName ...string) {
	if len(forceConfigName) > 0 {
		p.Cfg = config.GetConfig(forceConfigName[0])
	} else {
		p.Cfg = config.GetConfig()
	}
	/* The Gorilla router tests all the routes in order. The first one wins, the rest are ignored */
	p.Router = mux.NewRouter()

	/* The zap logger is fast and the sugared variant makes simple JSON as key-value pairs */
	p.Log = logger.GetLogger()

	// prefix for all calls to this proxy server e.g. /limelm
	appPrefix := p.Cfg.GetString("app_prefix")

	// make a subrouter to remove the prefix for all sub modules
	p.Log.Infow(p.Cfg.GetString("app_name")+": SubRoutePath ", "prefix", appPrefix)

	// `subrouter` provides paths with the prefix stripped off
	subrouter := p.Router.PathPrefix(appPrefix).Subrouter()

	//now add the limelm routes to the subrouter
	limelm.AddToRouter(subrouter)

	//static router - the backstop if no special router handles the request
	static.AddToRouter(p.Router)

	//print all the registered routes
	p.Log.Infow(p.Cfg.GetString("app_name")+": starting", "port", p.Cfg.GetInt("Port"))
	err := p.Router.Walk(p.gorillaWalkFn)
	if err != nil {
		p.Log.Fatal(err)
	}
}
