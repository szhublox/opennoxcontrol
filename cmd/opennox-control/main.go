package main

import (
	"log"
	"net/http"

	octl "github.com/szhublox/opennoxcontrol"
)

var (
	apiURL   = "http://127.0.0.1:18580"
	apiToken = "xyz"
)

const bindPort = "8080"

var (
	allowCmds = true
	bindHost  = "127.0.0.1:" + bindPort // set to 0.0.0.0 to allow external access
)

func main() {
	game := octl.NewGameHTTP(apiURL, apiToken)
	cp := octl.NewControlPanel(game, allowCmds)
	log.Fatal(http.ListenAndServe(bindHost, cp))
}
