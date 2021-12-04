package opennoxcontrol

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// NewGameHTTP creates a new Game client implemented over HTTP protocol.
func NewGameHTTP(apiURL string, token string) Game {
	return &httpClient{baseURL: apiURL, token: token}
}

type httpClient struct {
	baseURL string
	token   string
}

func (c *httpClient) GameInfo() (Info, error) {
	api_url := c.baseURL + "/api/v0/game/info"
	var info Info
	resp, err := http.Get(api_url)
	if err != nil {
		return info, errors.New("[opennoxcontrol]: couldn't get game data")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, errors.New("[opennoxcontrol]: couldn't parse game data")
	}
	json.Unmarshal(body, &info)
	return info, nil
}

func (c *httpClient) post(call string, data string) error {
	api_url := c.baseURL + "/api/v0/game/" + call
	body := strings.NewReader(data)

	req, err := http.NewRequest("POST", api_url, body)
	if err != nil {
		return errors.New("[opennoxcontrol]: couldn't generate POST request")
	}
	if c.token != "" {
		req.Header.Set("X-Token", c.token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("[opennoxcontrol]: couldn't send POST request")
	}
	defer resp.Body.Close()
	return nil
}

func (c *httpClient) ListMaps() ([]Map, error) {
	resp, err := http.Get(c.baseURL + "/api/v0/maps/")
	if err != nil {
		return nil, fmt.Errorf("[opennoxcontrol]: couldn't get map list: %v", err)
	}
	defer resp.Body.Close()
	var list []Map
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("[opennoxcontrol]: couldn't parse map list: %v", err)
	}
	return list, nil
}

func (c *httpClient) ChangeMap(name string) error {
	return c.post("map", name)
}

func (c *httpClient) Command(cmd string) error {
	return c.post("cmd", cmd)
}
