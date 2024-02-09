package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const HostUrl string = "https://gate.sendbird.com/api/v2"

type SendbirdClient struct {
	HostURL string
	ApiKey  string

	httpClient *http.Client
}

func New(host string, apiKey string) *SendbirdClient {
	c := SendbirdClient{
		HostURL:    HostUrl,
		ApiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	if host != "" {
		c.HostURL = host
	}

	return &c
}

func (c *SendbirdClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("SENDBIRDORGANIZATIONAPITOKEN", c.ApiKey)

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *SendbirdClient) buildUrl(path string, params ...string) (string, error) {
	url := c.HostURL + path

	if len(params) > 0 {
		fmt.Println(params)
	}

	return url, nil
}
