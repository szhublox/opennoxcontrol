package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

func get_info() (Info, error) {
	api_url := api_protocol + "://" + api_host + ":" + api_port +
		"/api/v0/game/info"
	var info Info

	resp, err := http.Get(api_url)
	if err != nil {
		return info, errors.New("[opennoxcontrol]: couldn't get game data")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, errors.New("[opennoxcontrol]: couldn't parse game data")
	}
	json.Unmarshal(body, &info)
	defer resp.Body.Close()
	return info, nil
}

func nox_curl_post(call string, data string) error {
	api_url := api_protocol + "://" + api_host + ":" + api_port +
		"/api/v0/game/" + call
	body := strings.NewReader(data)

	req, err := http.NewRequest("POST", api_url, body)
	if err != nil {
		return errors.New("[opennoxcontrol]: couldn't generate POST request")
	}
	req.Header.Set("X-Token", "xyz")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("[opennoxcontrol]: couldn't send POST request")
	}
	defer resp.Body.Close()
	return nil
}

func gameSetMap(name string) error {
	return nox_curl_post("map", name)
}

func gameCommand(cmd string) error {
	return nox_curl_post("cmd", cmd)
}
