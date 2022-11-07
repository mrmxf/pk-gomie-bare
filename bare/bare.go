package gomie

// Copyright (C) 2022 Dalet UK Ltd
// All rights reserved. Not for external use.

import (
	"fmt"
	"llmProxy/gomie/logger"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/** The GOlang MrMXF Internet Engine - gomie
 * a small and fast logging server for local deployment and proxies
 * based on the Gorilla golang server framework
 *
 *   - gomie.Setup  - default setup strategy (replace with your own)
 *   - gomie.Bare   - the naked server
 *   - gomie.Logger - a JSON oriented logger producing TIMESTAMP  LEVEL  SOURCE  MESSAGE {JSON}
 *   - gomie.Static - serve static files
 */

var appName string

// `router`     gorillaMux router that will allow CORS
// `port`       the int port to listen on
// `theAppName` for logging in a unified enviornment
// `wg`         a wait group for controlling shutdown
func (p *Proxy) Bare(router *mux.Router, port int, theAppName string, wg *sync.WaitGroup) *http.Server {

	appName = theAppName
	log := logger.GetLogger()

	//return a bare static server for the current folder
	//fire up the server
	location := fmt.Sprintf(":%v", port)
	log.Infow(appName+": listening", "state", "listening", "port", port)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a catch-all route"))
	})
	// loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// set up CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// start server listen
	// with error handling+/*

	// err := http.ListenAndServe(location, handlers.CORS(originsOk, headersOk, methodsOk)(loggedRouter))
	srv := &http.Server{
		Addr:    location,
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(router),
		// Good practice: enforce timeouts for servers!
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// always returns error. ErrServerClosed on graceful close
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		// unexpected error. port in use?
		log.Errorf("%v: ListenAndServe(port:%v): %v", appName, port, err)
		log.Sync()
	}

	// //set up our goroutine so that the server can be gracefuly sutdown if needed
	// go func() {
	// 	// let main know we are done cleaning up
	// 	defer wg.Done()

	// 	// always returns error. ErrServerClosed on graceful close
	// 	err := srv.ListenAndServe()
	// 	if err != http.ErrServerClosed {
	// 		// unexpected error. port in use?
	// 		log.Errorf("%v: ListenAndServe(port:%v): %v", appName, port, err)
	// 		log.Sync()
	// 	}
	// }()

	return srv
}
