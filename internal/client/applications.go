package client

import (
	"bytes"
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

func (c *SendbirdClient) CreateApplication(app_name string, region_key string) (*Application, error) {
	var application Application

	jsonBody := []byte(fmt.Sprintf(`{"app_name": "%s", "region_key": "%s"}`, app_name, region_key))
	readerBody := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/applications", c.HostURL), readerBody)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &application)
	if err != nil {
		return nil, err
	}
	return &application, err
}

func (c *SendbirdClient) DeleteApplication(appId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/applications/%s", c.HostURL, appId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}
