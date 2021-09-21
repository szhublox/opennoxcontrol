package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var api_protocol = "http"
var api_host = "127.0.0.1"
var api_port = "18580"

type Player struct {
	Name  string `json:"name"`
	Class string `json:"class"`
}

type PlayerInfo struct {
	Cur  int      `json:"cur"`
	Max  int      `json:"max"`
	List []Player `json:"list"`
}

type Info struct {
	Name       string     `json:"name"`
	Map        string     `json:"map"`
	Mode       string     `json:"mode"`
	Vers       string     `json:"vers"`
	PlayerInfo PlayerInfo `json:"players"`
}

func get_info() Info {
	api_url := api_protocol + "://" + api_host + ":" + api_port +
		"/api/v0/game/info"
	var info Info

	resp, err := http.Get(api_url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(body, &info)
	defer resp.Body.Close()
	return info
}

func nox_curl_post(call string, data string) {
	api_url := api_protocol + "://" + api_host + ":" + api_port +
		"/api/v0/game/" + call
	body := strings.NewReader(data)

	req, err := http.NewRequest("POST", api_url, body)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("X-Token", "xyz")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
}
