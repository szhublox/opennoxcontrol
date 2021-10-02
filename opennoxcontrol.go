package main

import (
	"log"
	"net/http"
)

var api_protocol = "http"
var api_host = "127.0.0.1"
var api_port = "18580"
var api_token = "xyz"

var bind_local = true // 127.0.0.1 instead of 0.0.0.0
var bind_port = "8080"

func main() {
	var bind_host string

	if bind_local {
		bind_host = "127.0.0.1:" + bind_port
	} else {
		bind_host = "0.0.0.0:" + bind_port
	}
	game := NewGameHTTP(api_protocol, api_host, api_port, api_token)
	cp := NewControlPanel(game, bind_local)
	log.Fatal(http.ListenAndServe(bind_host, cp))
}
