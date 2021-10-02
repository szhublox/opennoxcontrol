package main

import (
	"flag"
	"log"
	"net/http"

	octl "github.com/szhublox/opennoxcontrol"
)

var (
	apiURL    = flag.String("api", "http://127.0.0.1:18580", "base URL for the API endpoint")
	apiToken  = flag.String("token", "", "API token for controlling the server")
	allowCmds = flag.Bool("cmds", true, "allow sending commands to the server")
	bindHost  = flag.String("host", "127.0.0.1:8080", "host to listen on (:8080 to allow external access)")
)

func main() {
	flag.Parse()
	game := octl.NewGameHTTP(*apiURL, *apiToken)
	cp := octl.NewControlPanel(game, *allowCmds)
	log.Fatal(http.ListenAndServe(*bindHost, cp))
}
