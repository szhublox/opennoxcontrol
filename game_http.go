package opennoxcontrol

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, errors.New("[opennoxcontrol]: couldn't parse game data")
	}
	json.Unmarshal(body, &info)
	defer resp.Body.Close()
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

func (c *httpClient) ChangeMap(name string) error {
	return c.post("map", name)
}

func (c *httpClient) Command(cmd string) error {
	return c.post("cmd", cmd)
}
