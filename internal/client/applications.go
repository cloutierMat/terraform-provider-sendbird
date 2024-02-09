package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Region struct {
	Key  string `json:"region_key"`
	Name string `json:"region_name"`
}

type Application struct {
	Id        string `json:"app_id"`
	Name      string `json:"app_name"`
	ApiToken  string `json:"api_token"`
	CreatedAt string `json:"created_at"`
	Region    Region `json:"region"`
}

func (c *SendbirdClient) GetApplication(appId string) (*Application, error) {
	var application Application

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/applications/%s", c.HostURL, appId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &application)
	if err != nil {
		return nil, err
	}
	return &application, nil
}
