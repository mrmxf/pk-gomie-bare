# pk-gomie-bare

**GO**lang **M**rmxf **I**nternet **E**ngine - The bare version with no plugins

## Why?

I found myself making a number of custome servers for test & deployment
projects, much of which was the same code over & over again. Gomie-Bare is a
mini-development server that is easy to create & configure using a few lines of
code.

## Architecture

[GorillaMux](github.com/gorilla/mux) is used to handle routing,
[zap](go.uber.org/zap) for logging,
[viper](github.com/spf13/viper) fore layered config
[godotenv](github.com/joho/godotenv) for local secrets
and
[convey]9github.com/smartystreets/goconvey for testing

an embedded yaml config file gives the base configuration and other configs can
be layered on top (see below)

## Basic usage

1. Go get the package

```
go get github.com/mrmxf/pk-gomie-bare
```

Create a new golang project

```
go mod init myserver
```
 create `myserver.go`

 ```
package main

import (
	"github.com/mrmxf/pk-gomie-bare"
    "sync"
)

// Main Entry Point
func main() {
	gb := &dltproxy.Proxy{}
	gb.Setup()

	//setup the configuration
	cfg := cb.Cfg

    gb.Log()


	//create a waitgroup in case we need to shut down the server gracefully
	gbWg := &sync.WaitGroup{}
	gbWg.Add(1)

	//start the proxy server without handling graceful shutdown
	_ = gb.Bare(gb.Router, cfg.GetInt("Port"), cfg.GetString("Name"))
}

 ```